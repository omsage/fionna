package android

import (
	"fionna/android/gadb"
	"github.com/gorilla/websocket"
	"os/exec"
)

var (
	client gadb.Client
)

func init() {
	var err error
	client, err = gadb.NewClient()
	if err != nil {
		cmd := exec.Command("adb", "start-server")
		cmd.Run()
		var err1 error
		client, err1 = gadb.NewClient()
		if err1 != nil {
			panic("failed to connect adb server")
		}
	}
}

var upGrader websocket.Upgrader

func Init(upgrader websocket.Upgrader) {
	upGrader = upgrader
}
