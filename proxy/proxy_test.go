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

package proxy_test

import (
	"net/http"
	"testing"

	"gitlab.cern.ch/fts/lmt/mock"
	"gitlab.cern.ch/fts/lmt/proxy"
	"golang.org/x/net/websocket"
)

func init() {
	http.Handle("/socket", websocket.Handler(proxy.Handler))
	go http.ListenAndServe(":8081", nil)
}

const (
	origin   = "http://client.mock"
	server   = "ws://localhost:8081/socket"
	filepath = "sample.txt"
)

func TestRegisterClient(t *testing.T) {
	t.Log("Given the need to test that a WebSocket client can register itself.")
	conf, err := websocket.NewConfig(server, origin)
	if err != nil {
		t.Fatal(err)
	}
	c := mock.NewClient(conf)
	mock.Register(c, filepath)
	<-proxy.WsReader(c)
	if _, ok := proxy.Clients[c.Ws.LocalAddr().String()]; !ok {
		t.Fatalf("No endpoint has been associated with client %s",
			c.Ws.LocalAddr().String())
	}
	t.Logf("Client %s has been associated with the endpoint %s\n",
		c.Ws.LocalAddr().String(), proxy.Clients[c.Ws.LocalAddr().String()])
}
