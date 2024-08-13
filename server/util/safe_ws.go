package util

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type SafeWebsocket struct {
	ws   *websocket.Conn
	lock sync.Mutex
}

func NewSafeWebsocket(ws *websocket.Conn) *SafeWebsocket {
	return &SafeWebsocket{ws: ws}
}

func (s *SafeWebsocket) ReadJSON(v interface{}) error {
	var err error
	//s.lock.Lock()
	err = s.ws.ReadJSON(v)
	//s.lock.Unlock()
	return err
}

func (s *SafeWebsocket) WriteJSON(v interface{}) error {
	var err error
	s.lock.Lock()
	if s.ws != nil {
		err = s.ws.WriteJSON(v)
	} else {
		var data []byte
		data, err = json.Marshal(&v)
		fmt.Println(string(data))
	}
	s.lock.Unlock()
	return err
}
