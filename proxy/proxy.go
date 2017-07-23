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
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Clients maps a WebSocket connection to an endpoint.
var Clients map[string]string

func init() {
	Clients = make(map[string]string)
}

// Waits for and accepts the next TCP connection to the listener.
func asyncAccept(listener net.Listener) (<-chan net.Conn, error) {
	cc := make(chan net.Conn)

	go func() {
		defer close(cc)
		conn, err := listener.Accept()
		if err == nil {
			log.Info("Received connection!")
			cc <- conn
		}
	}()

	return cc, nil
}

// Wire calls a proxy function to stream data from client to service.
func Wire(c *Client, conn net.Conn) {
	log.Info("Wiring ", c.Ws.RemoteAddr(), " <=> ", conn.RemoteAddr())
	HTTPProxy(c, conn)
	log.Info("Wiring terminated")
}

func RegisterClient(uuid, endpoint string) {
	Clients[uuid] = endpoint
}

// Listen assigns a TCP listener for the client and starts the proxy service.
func Listen(c *Client) {
	// Listen on a new TCP port
	listenAddress := net.TCPAddr{}
	listener, err := net.ListenTCP("tcp", &listenAddress)
	if err != nil {
		log.Error(err)
		return
	}
	defer listener.Close()
	log.Info("Listening on ", listener.Addr().String())

	// Register client.
	// TODO: Generate a UUID for each client.
	Clients[c.Ws.RemoteAddr().String()] = listener.Addr().String()
	log.Infof("Client %s has been associated with the endpoint %s\n",
		c.Ws.RemoteAddr().String(), listener.Addr().String())

	// Recieve the onopen message from client
	var m CtrlMsg
	websocket.JSON.Receive(c.Ws, &m)

	closed := pingWebSocket(c)
	for {
		websocket.JSON.Send(c.Ws, []byte(fmt.Sprint("LISTEN ", listener.Addr())))
		cc, err := asyncAccept(listener)
		if err != nil {
			log.Error(err)
			break
		}

		var conn net.Conn
		var ok bool

		// wait until either a request has been recieved or the websocket has
		// been closed
		select {
		case conn, ok = <-cc:
			if !ok {
				log.Warn("Listener failed")
				return
			}
			log.Info("Received request")

		case closed := <-closed:
			log.Warn(closed)
			return
		}

		Wire(c, conn)
		conn.Close()
	}

	log.Info("Proxy finished")
}

// Handler is the WebSocket handler
func Handler(ws *websocket.Conn) {
	defer ws.Close()
	log.Info("Websocket proxy initiated")
	c := NewClient(ws)
	Listen(c)
	log.Info("Done here")
}
