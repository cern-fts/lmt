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
	log "github.com/Sirupsen/logrus"
	"gitlab.cern.ch/fts/lmt/lmt"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"sync"
	"time"
)

// Waits for and accepts the next TCP connection to the listener.
func asyncAccept(listener net.Listener) (<-chan net.Conn, error) {
	cc := make(chan net.Conn)

	go func() {
		defer close(cc)
		conn, err := listener.Accept()
		if err == nil {
			log.Info("Received connection!")
			cc <- conn
		}
	}()

	return cc, nil
}

// Calls io.Copy on a WaitGroup.
func dump(dst io.Writer, src io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	io.Copy(dst, src)
}

// Starts two seperate go routines to pipe data through the proxy.
func Pipeline(ws *websocket.Conn, conn net.Conn) {
	log.Info("Wiring ", ws.RemoteAddr(), " <=> ", conn.RemoteAddr())
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go dump(ws, conn, wg)
	go dump(conn, ws, wg)
	wg.Wait()
	log.Info("Wiring terminated")
}

// Continuously pings the WebSocket to make sure it is still open. If it fails
// to write to the WebSocket (client disconnected), it returns an error.
func pingWebSocket(ws *websocket.Conn) <-chan error {
	closed := make(chan error)
	go func() {
		for {
			if _, err := ws.Write([]byte("PING")); err != nil {
				closed <- err
				close(closed)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return closed
}

// Assigns a TCP listener for each websocket and maintains the connections.
func ProxyListen(ws *websocket.Conn) {
	defer ws.Close()
	listenAddress := net.TCPAddr{}

	listener, err := net.ListenTCP("tcp", &listenAddress)
	if err != nil {
		log.Error(err)
		return
	}
	transfer := lmt.Transfer{
		Origin:       ws.RemoteAddr().String(),
		ListenerAddr: listener.Addr().String(),
	}
	log.Info("Listening on ", transfer.ListenerAddr)
	defer listener.Close()

	closed := pingWebSocket(ws)

	for {
		ws.Write([]byte(fmt.Sprint("LISTEN ", listener.Addr())))

		cc, err := asyncAccept(listener)
		if err != nil {
			log.Error(err)
			break
		}

		var conn net.Conn
		var ok bool

		select {
		case conn, ok = <-cc:
			if !ok {
				log.Warn("Listener failed")
				return
			}
		case closed := <-closed:
			log.Warn(closed)
			return
		}

		Pipeline(ws, conn)
		conn.Close()
	}

	log.Info("Proxy finished")
}

// WebSocket Handler
func WebSocketProxy(ws *websocket.Conn) {
	log.Info("Websocket proxy initiated")
	ProxyListen(ws)
	log.Info("Done here")
}
