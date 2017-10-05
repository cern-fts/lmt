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

// transfer represents a transfer request submitted by a client.
type transfer struct {
	client   *client
	fileData *fileData
	identity string
	endPoint string
}

// fileData represents the details of the file to be transfered.
type fileData struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

// ctrlMsg is used to exchange messages with the client over a WebSocket
// connection.
type ctrlMsg struct {
	Action string `json:"action,omitempty"`
	Data   string `json:"data,omitempty"`
}

var (
	pingMsg = ctrlMsg{
		Action: "ping",
	}
	readyMsg = ctrlMsg{
		Action: "ready",
	}
	finishedMsg = ctrlMsg{
		Action: "finished",
	}
)
