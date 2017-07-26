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

	"github.com/Sirupsen/logrus"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func responseHeaders(t *transfer) http.Header {
	return http.Header{
		"Content-Type":   {"application/octet-stream"},
		"Accept-Ranges":  {"bytes"},
		"Content-Length": {string(t.FileData.Size)},
	}
}

// httpProxy handles HTTP requests from service.
func httpProxy(c *client, conn net.Conn) {
	// Read and parse incoming HTTP request
	buf := bufio.NewReader(conn)
	req, err := http.ReadRequest(buf)
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": "http_request_parse_error",
			"data":  err,
		}).Fatal(err)
	}
	// Log request
	dumpReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Error(err)
	}
	log.WithFields(logrus.Fields{
		"event": "http_request",
		"data":  string(dumpReq),
	}).Info(string(dumpReq))

	switch req.Method {
	case "HEAD":

		r := http.Response{
			Status:     "200",
			StatusCode: 200,
			Header:     responseHeaders(Clients[c.ID]),
		}
		r.Write(conn)

	case "GET":

		err := c.sendMsg(&readyMsg)
		if err != nil {
			log.WithFields(logrus.Fields{
				"event": "client_communication_error",
				"data":  err,
			}).Fatal(err)
		}
		frame := c.readFrame()
		r := http.Response{
			Status:     "200",
			StatusCode: 200,
			Header:     responseHeaders(Clients[c.ID]),
			Body:       nopCloser{bytes.NewBuffer(frame)},
		}
		r.Write(conn)
	}
}
