package mock

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"io/ioutil"
)

// ClientMock mocks a browser.
type ClientMock struct {
	Config *websocket.Config
	ws     *websocket.Conn
}

// Connect connects a new client to a web socket.
func (c *ClientMock) Connect() (*websocket.Conn, error) {
	log.Info("Connecting to ", c.Config.Location)
	ws, err := websocket.DialConfig(c.Config)
	if err != nil {
		log.Fatal(err)
	}
	c.ws = ws
	return ws, err
}

// UploadFile reads the whole file named by filepath and writes its content
// to the websocket.
func (c *ClientMock) UploadFile(filepath string) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.ws.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
