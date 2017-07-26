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
	"net"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Clients maps a WebSocket connection to an endpoint.
var Clients map[string]*transfer

func init() {
	Clients = make(map[string]*transfer)
}

// asyncAccept waits for and accepts the next TCP connection to the listener.
func asyncAccept(listener net.Listener) (<-chan net.Conn, error) {
	cc := make(chan net.Conn)

	go func() {
		defer close(cc)
		conn, err := listener.Accept()
		if err == nil {
			cc <- conn
		}
	}()

	return cc, nil
}

// tunnel calls a proxy function to stream data from client to service.
func tunnel(c *client, conn net.Conn) {
	log.WithFields(logrus.Fields{
		"event": "wiring_started",
	}).Info("Wiring ", c.Ws.RemoteAddr(), " <=> ", conn.RemoteAddr())
	httpProxy(c, conn)
	log.WithFields(logrus.Fields{
		"event": "wiring_finished",
	}).Info("Wiring finished")

}

// registerClient registers a new transfer.
func registerClient(uuid, origin, endpoint string, f *FileData) {
	Clients[uuid] = &transfer{
		Origin:   origin,
		Endpoint: endpoint,
		FileData: f,
	}
}

// listen assigns a TCP listener for the client and starts the proxy service.
func listen(c *client) {
	// Listen on a new TCP port
	listenAddress := net.TCPAddr{}
	listener, err := net.ListenTCP("tcp", &listenAddress)
	if err != nil {
		log.Error(err)
		return
	}
	defer listener.Close()
	log.WithFields(logrus.Fields{
		"event": "new_listener",
		"data":  listener.Addr().String(),
	}).Info("Listening on ", listener.Addr().String())

	// Recieve the onopen message from client
	var f FileData
	err = websocket.JSON.Receive(c.Ws, &f)
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": "onopen_message_error",
		}).Error(err)
	}
	log.WithFields(logrus.Fields{
		"event": "onopen_message_success",
		"data":  f,
	}).Info("Recieved JSON from websocket")

	// Register client.
	registerClient(c.ID, c.Ws.RemoteAddr().String(), listener.Addr().String(), &f)
	log.WithFields(logrus.Fields{
		"event": "client_registered",
		"data":  listener.Addr().String(),
	}).Infof("Client %s has been associated with the endpoint %s",
		c.Ws.RemoteAddr().String(), listener.Addr().String())

	endpointMsg := &ctrlMsg{
		Action: "transfer",
		Data:   listener.Addr().String(),
	}
	c.sendMsg(endpointMsg)
	closed := c.ping()
	for {
		// websocket.JSON.Send(c.Ws, []byte(fmt.Sprint("LISTEN ", listener.Addr())))
		cc, err := asyncAccept(listener)
		if err != nil {
			log.Error(err)
			break
		}

		log.WithFields(logrus.Fields{
			"event": "recieved_tcp_connection",
		}).Info("Recieved TCP connection")

		var conn net.Conn
		var ok bool
		// wait until either a request has been recieved or the websocket has
		// been closed
		select {
		case conn, ok = <-cc:
			if !ok {
				log.WithFields(logrus.Fields{
					"event": "listener_failed",
				}).Warn("Listener failed")
				return
			}
			log.WithFields(logrus.Fields{
				"event": "recieved_request",
			}).Info("Received request")

		case closed := <-closed:
			log.WithFields(logrus.Fields{
				"event": "websocket_closed",
			}).Warn(closed)
			return
		}

		tunnel(c, conn)
		conn.Close()
	}

	log.WithFields(logrus.Fields{
		"event": "listener_finised",
	}).Info("listener finished")
}

// Handler is the WebSocket handler
func Handler(ws *websocket.Conn) {
	defer ws.Close()
	log.WithFields(logrus.Fields{
		"event": "websocket_handler_initiated",
		"data":  ws.RemoteAddr().String(),
	}).Info("Websocket handler initiated")
	// Assign a UUID for the WebSocket connection
	c := newClient(ws)
	// Listen for service requests.
	listen(c)

	log.WithFields(logrus.Fields{
		"event": "websocket_handler_finished",
		"data":  ws.RemoteAddr().String(),
	}).Info("Websocket proxy finished")
}
