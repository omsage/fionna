package scrcpy_client

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"fionna/android/gadb"
	"fionna/entity"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	deviceServerPath = "/data/local/tmp/scrcpy-server.jar"
	sockName         = "scrcpy"
)

var (
	// 来源：https://github.com/aoliaoaoaojiao/scrcpy
	//go:embed scrcpy-server
	scrcpyBytes    []byte
	scrcpyPtrLen   int32
	scrcpyPtrLock  sync.Mutex
	scrcpyPtrStore []*Scrcpy
)

func StoreScrcpyPtr(v *Scrcpy) int32 {

	id := atomic.AddInt32(&scrcpyPtrLen, 1) - 1
	scrcpyPtrLock.Lock()
	defer scrcpyPtrLock.Unlock()
	scrcpyPtrStore = append(scrcpyPtrStore, v)
	log.Debugf("store scrcpy object,the id is: %d", id)
	return id
}

func RestoreScrcpyPtr(ptr int32) *Scrcpy {
	scrcpyPtrLock.Lock()
	defer scrcpyPtrLock.Unlock()
	log.Debugf("get scrcpy object form id %d", ptr)
	return scrcpyPtrStore[ptr]
}

type Scrcpy struct {
	dev        *gadb.Device
	localPort  int
	isSendH246 bool
	scrcpyLn   net.Listener
	//scrcpyControl    *Control
	videoSocket net.Conn
	//controlSocket    net.Conn
	sizeInfoSocket   net.Conn
	lock             sync.Mutex
	forwardWs        *websocket.Conn
	exitCallBackFunc context.CancelFunc
	exitCtx          context.Context
}

func NewScrcpy(device *gadb.Device, ctx context.Context, forwardWs *websocket.Conn) *Scrcpy {
	// todo
	ln, err := net.Listen("tcp", ":0") // 0表示随机端口
	if err != nil {
		return nil
	}

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		return nil
	}

	ctx, exitFunc := context.WithCancel(ctx)

	log.Debugf("create scrcpy object,service port %d", tcpAddr.Port)

	return &Scrcpy{
		dev:              device,
		scrcpyLn:         ln,
		localPort:        tcpAddr.Port,
		exitCtx:          ctx,
		exitCallBackFunc: exitFunc,
		forwardWs:        forwardWs,
	}
}

func (s *Scrcpy) Start(pic entity.ScrcpyPic) {
	// todo
	var err error
	err = s.dev.Push(bytes.NewReader(scrcpyBytes), deviceServerPath, time.Now())
	if err != nil {
		panic(err)
	}

	err = s.dev.ReverseLocalAbstract(sockName, s.localPort)
	if err != nil {
		panic(err)
	}

	go func() {
		<-s.exitCtx.Done()
		s.clientStop()
	}()

	s.startServer()

	s.runBinary(pic)

}

func (s *Scrcpy) runBinary(pic entity.ScrcpyPic) {
	var output io.Reader

	var maxSize int

	switch pic {
	case entity.ScrcpyPicLow:
		maxSize = 640
	case entity.ScrcpyPicMid:
		maxSize = 1280
	case entity.ScrcpyPicHeight:
		maxSize = 1920
	}

	output, err := s.dev.RunShellLoopCommand(fmt.Sprintf("CLASSPATH=/data/local/tmp/scrcpy-server.jar app_process / com.genymobile.scrcpy.Server v2.2  log_level=debug max_size=0 max_fps=60 control=false max_size=%d audio=false audio=false size_info=true", maxSize))
	if err != nil {
		log.Error("execute scrcpy err:", err)
		s.exitCallBackFunc()
		panic(err)
	}
	var isRelease sync.WaitGroup

	isRelease.Add(1)

	// 运行scrcpy
	go func() {
		var byteDatas = make([]byte, 1024)
		n, err := output.Read(byteDatas)
		if err != nil {
			log.Error("start scrcpy err:", err)
			s.exitCallBackFunc()
			return
		}
		if !strings.Contains(string(byteDatas[:n]), "Device") {
			log.Error("not start scrcpy:", string(byteDatas[:n]))
			s.exitCallBackFunc()
			return
		}
		isRelease.Done()

		for {
			select {
			case <-s.exitCtx.Done():
				return
			default:
				n, err = output.Read(byteDatas)
				if err != nil {
					// 如果发生超时错误，你可以根据具体情况进行处理
					if err, ok := err.(net.Error); ok && err.Timeout() {
						time.Sleep(1 * time.Second)
						continue
					}
					log.Error("get scrcpy binary out:", err)
					s.exitCallBackFunc()
					return
				}
				log.Debug(string(byteDatas[:n]))
			}
		}
	}()
	log.Debugf("start scrcpy server!")
	isRelease.Wait()
}

func (s *Scrcpy) startServer() {
	go func() {
		var err error
		s.videoSocket, err = s.scrcpyLn.Accept()
		if err != nil {
			log.Error("get scrcpy video socket err,", err)
			return
		}
		// 解析和转发video socket
		go func() {
			s.videoParse()
		}()

		//s.controlSocket, err = s.scrcpyLn.Accept()
		//s.scrcpyControl = NewControl(s.controlSocket)
		if err != nil {
			log.Error("get scrcpy control socket err,", err)
			return
		}
		// 解析和转发sizeinfo socket
		s.sizeInfoSocket, err = s.scrcpyLn.Accept()
		if err != nil {
			log.Error("get scrcpy rotation socket err,", err)
			return
		}
		go func() {
			s.sizeInfoParse()
		}()
	}()

}

func (s *Scrcpy) videoParse() {
	buffer := make([]byte, 64)
	_, err := s.videoSocket.Read(buffer)
	if err != nil {
		log.Error("get scrcpy device info err,", err)
		s.exitCallBackFunc()
	}
	buffer = make([]byte, 12)
	_, err = s.videoSocket.Read(buffer)
	if err != nil {
		log.Error("get scrcpy device width and height info err,", err)
		s.exitCallBackFunc()
	}
	//binary.BigEndian.Uint32(buffer[:4])
	//s.width = int(binary.BigEndian.Uint32(buffer[4:8]))
	//s.height = int(binary.BigEndian.Uint32(buffer[8:12]))

	go func() {
		s.writeH264()
	}()
}

func (s *Scrcpy) sizeInfoParse() {
	for {
		select {
		case <-s.exitCtx.Done():
			return
		default:
			lengthBuffer := make([]byte, 4)
			n, err1 := s.sizeInfoSocket.Read(lengthBuffer)
			if err1 != nil {
				log.Error("get scrcpy size info length err,", err1)
				s.exitCallBackFunc()
				break
			}
			if n == 0 {
				continue
			}

			dataLen := binary.BigEndian.Uint32(lengthBuffer)

			sizeInfoDataBuffer := make([]byte, dataLen-4)

			_, err1 = s.sizeInfoSocket.Read(sizeInfoDataBuffer)
			if err1 != nil {
				log.Error("get scrcpy size info err,", err1)
				s.exitCallBackFunc()
				break
			}

			sizeInfoData := &entity.ScrcpySizeInfo{}

			err1 = json.Unmarshal(sizeInfoDataBuffer, sizeInfoData)
			if err1 != nil {
				panic(err1)
			}

			if s.forwardWs != nil {
				s.lock.Lock()
				err := s.forwardWs.WriteJSON(
					entity.ScrcpySizeInfoMessage{
						MessageType: entity.ScrcpySizeInfoType,
						Data:        sizeInfoData,
					},
				)
				s.lock.Unlock()
				if err != nil {
					log.Error("froward device rotation err,", err)
					s.exitCallBackFunc()
					break
				}
				s.isSendH246 = true
			}
		}

	}
}

func (s *Scrcpy) writeH264() {
	buf := make([]byte, 1024*1024*10)
	//file, err := os.OpenFile("data.h246", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()

	for {
		select {
		case <-s.exitCtx.Done():
			return
		default:
			if !s.isSendH246 {
				continue
			}
			bytesLen, err := s.videoSocket.Read(buf)
			if err != nil {
				// 如果发生超时错误，你可以根据具体情况进行处理
				if err1, ok := err.(net.Error); ok && err1.Timeout() {
					time.Sleep(1 * time.Second)
					continue
				} else {
					log.Error("get scrcpy h264 data err,", err)
					s.exitCallBackFunc()
					break
				}
			}

			//_, err = file.Write(buf[:bytesLen])
			//if err != nil {
			//	panic(err)
			//}
			//var lastNaluIndex = 0
			s.lock.Lock()
			err = s.forwardWs.WriteMessage(websocket.BinaryMessage, buf[:bytesLen])
			if err != nil {
				log.Error("forward h264 forward err,", err)
				s.exitCallBackFunc()
				break
			}
			s.lock.Unlock()
		}
	}
}

func (s *Scrcpy) ClientStop() {
	s.exitCallBackFunc()
}

func (s *Scrcpy) clientStop() {
	log.Debugf("close scrcpy server")
	if s.scrcpyLn != nil {
		s.dev.ReverseKillLocalAbstract(sockName)
		s.scrcpyLn.Close()
		s.scrcpyLn = nil
	}
	if s.videoSocket != nil {
		s.videoSocket.Close()
		s.videoSocket = nil
		//s.frameStop()
	}
	//if s.controlSocket != nil {
	//	s.controlSocket.Close()
	//	s.controlSocket = nil
	//}
	if s.forwardWs != nil {
		s.forwardWs.Close()
		s.forwardWs = nil
	}
}

//func (s *Scrcpy) Touch(touch *entity.ScrcpyTouch, touchID int) error {
//	return s.scrcpyControl.Touch(touch, touchID)
//}

func (s *Scrcpy) GetDevice() *gadb.Device {
	return s.dev
}
