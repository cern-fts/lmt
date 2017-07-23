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

package mock

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gitlab.cern.ch/fts/lmt/proxy"
	"golang.org/x/net/websocket"
)

type Msg struct {
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}

// NewClient creates a new WebSocket client connection.
func NewClient(c *websocket.Config) *proxy.Client {
	ws, err := websocket.DialConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	return proxy.NewClient(ws)
}

// Register registers a transfer request.
func Register(c *proxy.Client, filepath string) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	openMsg := Msg{
		Name: filepath,
		Size: len(contents),
	}
	websocket.JSON.Send(c.Ws, openMsg)
}
