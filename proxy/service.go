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
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// ServiceHandler handles HTTP requests from service (FTS) to the proxy.
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Log request
	dumpReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Error(err)
	}
	log.WithFields(logrus.Fields{
		"event": "http_request",
		"data":  string(dumpReq),
	}).Info(string(dumpReq))

	// Parse request vars and get transferID
	vars := mux.Vars(r)
	transferID := vars["id"]

	if transfer, found := Clients[transferID]; !found {
		// Transfer not found.
		w.WriteHeader(http.StatusNotFound)
	} else {
		// Transfer exists, set response headers.
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", transfer.fileData.Size))

		if r.Method == "GET" {
			c := transfer.client
			log.WithFields(logrus.Fields{
				"event": "wiring_started",
			}).Info("Wiring ", c.Ws.RemoteAddr(), " <=> ", r.RemoteAddr)
			// Notify client that the service is ready to start the transfer.
			err := c.sendMsg(&readyMsg)
			if err != nil {
				log.WithFields(logrus.Fields{
					"event": "client_communication_error",
					"data":  err,
				}).Error(err)
			}
			// Stream data from client (websocket connection) to service.
			io.Copy(w, c.Ws)
			log.WithFields(logrus.Fields{
				"event": "wiring_finished",
			}).Info("Wiring finished")
		}
	}
}
