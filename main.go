package main

import (
	log "github.com/Sirupsen/logrus"
	"gitlab.cern.ch/fts/lmt/mock"
	"gitlab.cern.ch/fts/lmt/proxy"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"path"
	"time"
)

var filepath string

func init() {
	cwd, _ := os.Getwd()
	filepath = path.Join(cwd, "mock", "sample.txt")
}

func main() {
	// origin & target URLs for websocket config
	var origin = "http://localhost"
	var target = "ws://localhost:8080/socket"
	http.Handle("/socket", websocket.Handler(proxy.WebSocketProxy))
	log.Info("Listening on http://localhost:8080")
	go func() {
		time.Sleep(time.Second * 1)
		config, err := websocket.NewConfig(target, origin)
		if err != nil {
			log.Fatal(err)
		}
		c := &mock.ClientMock{
			Config: config,
		}
		c.Connect()
		time.Sleep(time.Second * 1)
		c.UploadFile(filepath)
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
