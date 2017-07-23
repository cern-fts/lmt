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
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"

	log "github.com/Sirupsen/logrus"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// HTTPProxy handles HTTP requests from service.
func HTTPProxy(c *Client, conn net.Conn) {
	// Read HTTP request into buffer
	buf := bufio.NewReader(conn)
	req, err := http.ReadRequest(buf)
	if err != nil {
		panic(err)
	}
	// Log request
	dumpReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(dumpReq))

	// Respond to request.
	// TODO: write filedata to the response's body.
	switch req.Method {
	case "HEAD":
		r := http.Response{
			Status:     "200",
			StatusCode: 200,
		}
		r.Write(conn)

	case "GET":
		_ = c.SendMsg(&CtrlMsg{Action: "ready"})
		wire := WsReader(c)
		r := http.Response{
			Status:     "200",
			StatusCode: 200,
			Body:       nopCloser{bytes.NewBuffer(<-wire)},
		}
		r.Write(conn)
	}
}
