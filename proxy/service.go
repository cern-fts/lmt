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
	voms "gitlab.cern.ch/flutter/go-proxy"
)

// Config struct is used to unmarshal the configuration found in
// the config.yml file
type Config struct {
	Headers map[string]string `yaml:"additional_http_headers"`
}

// ResponseHeaders stores key-value pairs of additional HTTP headers
// that are configured via the config.yml file.
var ResponseHeaders map[string]string

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

	identity, err := X509Identity(r)
	if err != nil {
		log.Error(err)
	}
	// Parse request vars and get transferID.
	vars := mux.Vars(r)
	delegationID := vars["delegationID"]
	filename := vars["filename"]
	transferID := TransferID(delegationID, filename)

	if transfer, found := Transfers[transferID]; !found {
		// Transfer not found.
		w.WriteHeader(http.StatusNotFound)
	} else {
		// Transfer exists, check permissions
		if !CheckIdentity(transferID, voms.NameRepr(&identity)) {
			// FTS does not have permission to access the file.
			w.WriteHeader(http.StatusForbidden)
			log.WithFields(logrus.Fields{
				"event": "access_forbidden",
				"data":  string(voms.NameRepr(&identity)),
			}).Error(errAccessForbidden)
		} else {
			// FTS has the correct permissions.
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", transfer.fileData.Size))

			switch r.Method {
			case "HEAD":
				w.WriteHeader(http.StatusOK)
			case "GET":
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
				_, err = io.CopyN(w, c.Ws, transfer.fileData.Size)
				if err != nil {
					log.WithFields(logrus.Fields{
						"event": "proxy_tunneling_error",
						"data":  err,
					}).Error(err)
				}
				log.WithFields(logrus.Fields{
					"event": "wiring_finished",
				}).Info("Wiring finished")
				// Notify client that the transfer has been successfully completed.
				err = c.sendMsg(&finishedMsg)
				if err != nil {
					log.WithFields(logrus.Fields{
						"event": "client_communication_error",
						"data":  err,
					}).Error(err)
				}
				// Remove the transfer from the Transfers map and close the
				// websocket connection.
				err = c.close()
				if err != nil {
					log.WithFields(logrus.Fields{
						"event": "websocket_close_error",
						"data":  err,
					}).Error(err)
				}
			default:
				// Return 405 for all other request methods.
				// This is especially important when a "COPY" request is received,
				// so that FTS would later do a GET request and pull the data.
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}
}
