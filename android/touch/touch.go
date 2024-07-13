package touch

import (
	"bytes"
	_ "embed"
	"fionna/android/gadb"
	"fionna/entity"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	// source:https://github.com/aoliaoaoaojiao/AndroidTouch
	//go:embed lib/AndroidTouch.jar
	touchBytes []byte
)

const touchToolPath = "/data/local/tmp/atouch.jar"

type Touch struct {
	shellLoopConn net.Conn
}

func NewTouch(device *gadb.Device) *Touch {
	var err error
	err = device.Push(bytes.NewReader(touchBytes), touchToolPath, time.Now())
	if err != nil {
		panic(err)
	}
	conn, err := device.RunShellLoopCommandSock(fmt.Sprintf(
		"CLASSPATH=%s app_process / com.aoliaoaojiao.AndroidTouch.Run v2.2",
		touchToolPath))
	if err != nil {
		panic(err)
	}
	var isRelease sync.WaitGroup

	isRelease.Add(1)
	go func() {
		var byteDatas = make([]byte, 1024)
		n, err := conn.Read(byteDatas)
		if err != nil {
			logrus.Error("start android touch err:", err)
			return
		}
		if !strings.Contains(string(byteDatas[:n]), "Device") {
			isRelease.Done()
			logrus.Error("not start android touch:", string(byteDatas[:n]))
			return
		}
		isRelease.Done()

	}()
	isRelease.Wait()
	return &Touch{shellLoopConn: conn}
}

func (t *Touch) Touch(touchInfo entity.TouchInfo) {
	var cmd string
	switch touchInfo.TouchType {
	case entity.TOUCH_DOWN, entity.TOUCH_MOVE:
		cmd = fmt.Sprintf(
			"%s %s %f %f %d\n",
			"airtest",
			touchInfo.TouchType,
			touchInfo.X,
			touchInfo.Y,
			touchInfo.FingerID)
	case entity.TOUCH_UP:
		cmd = fmt.Sprintf(
			"%s %s %d\n",
			"airtest",
			touchInfo.TouchType,
			touchInfo.FingerID)
	}
	_, err := t.shellLoopConn.Write([]byte(cmd))
	if err != nil {
		logrus.Error(err)
	}
}
