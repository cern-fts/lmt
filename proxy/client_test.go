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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

var serverAddr string
var once sync.Once

// startServer starts a web server and handles incoming web socket
// connections on /socket.
func startServer() {
	http.Handle("/socket", websocket.Handler(ClientHandler))
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Info("Test WebSocket server listening on ", serverAddr)
}

func TestWebsocketDial(t *testing.T) {
	// Start the test server.
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s/socket", serverAddr)
	origin := "http://localhost"
	_, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientHandler(t *testing.T) {
	// Start the test server.
	once.Do(startServer)
	// Test cases.
	tt := []struct {
		name   string
		origin string
		url    string
		file   fileData
	}{
		{
			name:   "transfer registered successfully",
			origin: "http://localhost",
			url:    fmt.Sprintf("ws://%s/socket", serverAddr),
			file: fileData{
				Name: "sample.txt",
				Size: 3875,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Dial websocket.
			ws, err := websocket.Dial(tc.url, "", tc.origin)
			if err != nil {
				t.Fatal(err)
			}
			// Register transfer request.
			err = websocket.JSON.Send(ws, tc.file)
			if err != nil {
				t.Fatal(err)
			}

			// Read the service response.
			var endpointMsg ctrlMsg
			err = websocket.JSON.Receive(ws, &endpointMsg)
			if err != nil {
				t.Fatal(err)
			}
			log.Infof("%v", endpointMsg)
		})
	}
}
