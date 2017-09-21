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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitlab.cern.ch/fts/lmt/proxy"
	"golang.org/x/net/websocket"
)

var staticDir string

func init() {
	cwd, _ := os.Getwd()
	staticDir = path.Join(cwd, "static")
}

func main() {
	r := mux.NewRouter()
	// Serve static files.
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir+"/"))))
	// Endpoint to be called by the transfer service (FTS).
	r.HandleFunc("/transfer/{id}", proxy.ServiceHandler)
	// Endpoint to be called by the client (web browser).
	r.Handle("/socket", websocket.Handler(proxy.ClientHandler))
	// Parse command-line flags to set the listening address for the proxy
	// service and the base URL for the transfer endpoints.
	port := flag.String("listen", ":8080", "port to listen on")
	flag.Parse()
	addr := fmt.Sprintf("http://localhost%s", *port)
	proxy.BaseURL = fmt.Sprintf("%s/%s/", addr, "transfer")
	// Start the web service.
	log.Infof("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(*port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /")

	index := path.Join(staticDir, "index.html")
	fd, err := os.Open(index)
	if err != nil {
		log.Warn(err)
		http.NotFound(w, r)
		return
	}
	defer fd.Close()

	body, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Error(err)
	}
}
