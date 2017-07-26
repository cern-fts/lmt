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
	"time"

	"golang.org/x/net/websocket"
)

// Client represents a WebSocket client.
type client struct {
	ID string
	Ws *websocket.Conn
}

// newClient constructs a new WebSocket Client.
func newClient(ws *websocket.Conn) *client {
	uuid, err := NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return &client{
		ID: uuid,
		Ws: ws,
	}
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
				closed <- err
				close(closed)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return closed
}
