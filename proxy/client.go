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

// CtrlMsg is used to exchange messages with the client over a WebSocket
// connection.
type CtrlMsg struct {
	Action string `json:"action,omitempty"`
}

// Client represents a WebSocket client.
type Client struct {
	Ws *websocket.Conn
}

// NewClient constructs a new WebSocket Client.
func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		Ws: ws,
	}
}

// SendMsg encodes a CtrlMsg and sends it to the client.
func (c *Client) SendMsg(m *CtrlMsg) <-chan error {
	err := make(chan error)
	go func() {
		defer close(err)
		err <- websocket.JSON.Send(c.Ws, *m)
	}()
	return err
}

// WsReader waits on a WebSocket connection until a frame can be read.
func WsReader(c *Client) <-chan []byte {
	var frame []byte
	wire := make(chan []byte, len(frame))
	go func() {
		for {
			if err := websocket.Message.Receive(c.Ws, &frame); err == nil {
				wire <- frame
			}
		}
	}()
	return wire
}

// Continuously pings the WebSocket to make sure it is still open. If it fails
// to write to the WebSocket (client disconnected), it returns an error.
func pingWebSocket(c *Client) <-chan error {
	closed := make(chan error)
	go func() {
		for {
			if err := websocket.JSON.Send(c.Ws, []byte("PING")); err != nil {
				closed <- err
				close(closed)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return closed
}
