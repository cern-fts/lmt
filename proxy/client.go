/*
 * Copyright (c) CERN 2017
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proxy

import (
	"net/http/httputil"
	"time"

	"github.com/Sirupsen/logrus"
	voms "gitlab.cern.ch/flutter/go-proxy"
	"golang.org/x/net/websocket"
)

// client represents a WebSocket client.
type client struct {
	ID string
	Ws *websocket.Conn
}

// registerClient creates a new client and adds it to the Transfers map.
func registerClient(ws *websocket.Conn, transferID string, f *fileData) *client {
	req := ws.Request()
	dumpReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Error(err)
	}
	log.WithFields(logrus.Fields{
		"event": "ws_http_request",
		"data":  string(dumpReq),
	}).Info(string(dumpReq))

	identity, err := X509Identity(req)
	if err != nil {
		log.Error(err)
	}
	c := &client{
		ID: transferID,
		Ws: ws,
	}
	// add new transfer to the map.
	Transfers[transferID] = &transfer{
		client:   c,
		fileData: f,
		identity: voms.NameRepr(&identity),
		endPoint: BaseURL + transferID,
	}
	return c
}

// sendMsg encodes a ctrlMsg and sends it to the client.
func (c *client) sendMsg(m *ctrlMsg) error {
	return websocket.JSON.Send(c.Ws, *m)
}

// readFrame waits on a WebSocket connection until a frame can be read.
func (c *client) readFrame() []byte {
	var frame []byte
	for {
		if err := websocket.Message.Receive(c.Ws, &frame); err == nil {
			return frame
		}
	}
}

// Continuously pings the WebSocket to make sure it is still open. If it fails
// to write to the WebSocket (client disconnected), it returns an error.
func (c *client) ping() <-chan error {
	closed := make(chan error)
	go func() {
		for {
			if err := websocket.JSON.Send(c.Ws, pingMsg); err != nil {
				log.WithFields(logrus.Fields{
					"event": "ping_client_fail",
				}).Info("Failed to ping client.")
				closed <- err
				close(closed)
				return
			}
			log.WithFields(logrus.Fields{
				"event": "ping_client_success",
			}).Info("Pinged client.")
			time.Sleep(10 * time.Second)
		}
	}()
	return closed
}

// close removes the transfer from the transfers map and closes the
// corresponding websocket connection.
func (c *client) close() error {
	delete(Transfers, c.ID)
	return c.Ws.Close()
}

// ClientHandler handles incoming websocket connections.
func ClientHandler(ws *websocket.Conn) {
	defer ws.Close()
	log.WithFields(logrus.Fields{
		"event": "websocket_handler_initiated",
		"data":  ws.RemoteAddr().String(),
	}).Info("Websocket handler initiated")
	// Recieve the onopen message from client
	var f fileData
	err := websocket.JSON.Receive(ws, &f)
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": "onopen_message_error",
		}).Error(err)
	}
	log.WithFields(logrus.Fields{
		"event": "onopen_message_success",
		"data":  f,
	}).Info("Recieved JSON from websocket")

	// Register a new client.
	uid, err := NewUUID()
	c := registerClient(ws, uid, &f)
	log.WithFields(logrus.Fields{
		"event": "client_registered",
		"data":  c.ID,
	}).Infof("Client %s has been associated with ID %s",
		c.Ws.RemoteAddr().String(), c.ID)

	// Send endpoint URL to client.
	endpointMsg := &ctrlMsg{
		Action: "transfer",
		Data:   Transfers[c.ID].endPoint,
	}
	c.sendMsg(endpointMsg)

	closed := c.ping()
	for {
		// wait until the websocket has been closed
		select {
		case closed := <-closed:
			log.WithFields(logrus.Fields{
				"event": "websocket_closed",
			}).Warn(closed)
			return
		}

	}

	log.WithFields(logrus.Fields{
		"event": "websocket_handler_finished",
		"data":  ws.RemoteAddr().String(),
	}).Info("Websocket proxy finished")
}
