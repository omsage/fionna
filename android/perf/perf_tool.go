package perf

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"fionna/android/gadb"
	"fionna/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"path"
	"strings"
	"sync"
	"time"
)

type PerfTool struct {
	dev               *gadb.Device
	localPort         int
	width, height     int
	perfToolLn        net.Listener
	frameSocket       net.Conn
	exitCallBackFunc  context.CancelFunc
	exitCtx           context.Context
	perfFrameDataChan chan *entity.SysFrameInfo
}

func NewPerfTool(device *gadb.Device, exitCtx context.Context) *PerfTool {
	// todo
	ln, err := net.Listen("tcp", ":0") // 0表示随机端口
	if err != nil {
		return nil
	}

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		return nil
	}

	ctx, exitFunc := context.WithCancel(exitCtx)

	return &PerfTool{
		dev:               device,
		perfToolLn:        ln,
		localPort:         tcpAddr.Port,
		exitCtx:           ctx,
		exitCallBackFunc:  exitFunc,
		perfFrameDataChan: make(chan *entity.SysFrameInfo),
	}
}

var (
	//go:embed lib/PerfTool.jar
	fraemToolByte []byte

	//go:embed lib/arm64-v8a/libPerfTool.so
	libArm64 []byte

	//go:embed lib/armeabi-v7a/libPerfTool.so
	libArm32 []byte

	//go:embed lib/x86_64/libPerfTool.so
	lib86_64 []byte

	//go:embed lib/x86/libPerfTool.so
	lib86 []byte
)

const perfToolPath = "/data/local/tmp/omsage-PerfTool.jar"

const perfToolLibRemotePath = "/data/local/tmp"

const perfSockNamePre = "PerfTool"

func (s *PerfTool) Init() {
	var err error
	err = s.dev.Push(bytes.NewReader(fraemToolByte), perfToolPath, time.Now())
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	// 生成随机整数
	randomInt := rand.Intn(100000) // 生成0到99之间的随机整数

	abi, _ := s.dev.RunShellCommand("getprop ro.product.cpu.abi")

	libPushPath := path.Join(perfToolLibRemotePath, "libPerfTool.so")

	if strings.Contains(abi, "arm64-v8a") {
		s.dev.Push(bytes.NewReader(libArm64), libPushPath, time.Now())
	} else if strings.Contains(abi, "armeabi-v7a") {
		s.dev.Push(bytes.NewReader(libArm32), libPushPath, time.Now())
	} else if strings.Contains(abi, "x86") {
		s.dev.Push(bytes.NewReader(lib86), libPushPath, time.Now())
	} else {
		s.dev.Push(bytes.NewReader(lib86_64), libPushPath, time.Now())
	}
	err = s.dev.ReverseLocalAbstract(perfSockNamePre+fmt.Sprintf("_%d", randomInt), s.localPort)
	if err != nil {
		panic(err)
	}

	go func() {
		<-s.exitCtx.Done()
		s.clientStop()
	}()

	s.startServer()

	s.runBinary(randomInt)

}

func (s *PerfTool) runBinary(cid int) {
	var output io.Reader

	output, err := s.dev.RunShellLoopCommand(fmt.Sprintf(
		"LD_LIBRARY_PATH=/system/lib64:/system_ext/lib64:%s "+
			"CLASSPATH=%s "+
			"app_process / com.omsage.PerfTool.Run 1.0 cid=%d",
		perfToolLibRemotePath, perfToolPath, cid))
	if err != nil {
		s.clientStop()
		panic(err)
	}
	var isRelease sync.WaitGroup

	isRelease.Add(1)

	go func() {
		var bytesOutput = make([]byte, 4090)
		n, err := output.Read(bytesOutput)
		if err != nil {
			s.exitCallBackFunc()
			log.Error(err)
			return
		}
		if !strings.Contains(string(bytesOutput[:n]), "Device") {
			s.clientStop()
			// todo
			log.Error("start fail! output: " + string(bytesOutput[:n]))
			return
		}
		isRelease.Done()

		for {
			select {
			case <-s.exitCtx.Done():
				return
			default:
				n, err = output.Read(bytesOutput)
				if err != nil {
					// 如果发生超时错误，你可以根据具体情况进行处理
					if err1, ok := err.(net.Error); ok && err1.Timeout() {
						time.Sleep(1 * time.Second)
						continue
					}
					s.exitCallBackFunc()
					log.Error("frame output err,", err)
					return
				}
				log.Debug(string(bytesOutput[:n]))
			}

		}
	}()
	isRelease.Wait()
}

func (s *PerfTool) startServer() {
	go func() {
		var err error
		s.frameSocket, err = s.perfToolLn.Accept()
		if err != nil {
			// todo
			s.frameSocket = nil
			s.exitCallBackFunc()
		}

		s.getFrameSteam()
	}()

}

func (s *PerfTool) getFrameSteam() {
	for {
		select {
		case <-s.exitCtx.Done():
			return
		default:
			lengthBuffer := make([]byte, 4)
			n, err1 := s.frameSocket.Read(lengthBuffer)
			if err1 != nil {
				// todo
				s.exitCallBackFunc()
				break
			}
			if n == 0 {
				continue
			}

			dataLen := binary.BigEndian.Uint32(lengthBuffer)

			perfDataBuffer := make([]byte, dataLen-4)

			_, err1 = s.frameSocket.Read(perfDataBuffer)
			if err1 != nil {
				// todo
				s.exitCallBackFunc()
				break
			}
			perfData := &entity.SysFrameInfo{}
			err1 = json.Unmarshal(perfDataBuffer, perfData)
			if err1 != nil {
				panic(err1)
			}
			s.perfFrameDataChan <- perfData
		}
	}
}

func (s *PerfTool) GetFrame(getBackCall func(frame *entity.SysFrameInfo, code entity.ServerCode)) {
	ticker := time.NewTicker(1 * time.Second)
	isNoFirst := false
	go func() {
		for {
			<-ticker.C
			if s.perfToolLn == nil {
				return
			}
			select {
			case perfData, ok := <-s.perfFrameDataChan:
				isNoFirst = true
				if ok && getBackCall != nil {
					getBackCall(perfData, entity.RequestSucceed)
				}
			default:
				if isNoFirst {
					getBackCall(&entity.SysFrameInfo{Timestamp: time.Now().UnixMilli()}, entity.RequestSucceed)
				}
			}
		}
	}()
}

func (s *PerfTool) clientStop() {
	log.Debugf("close frame tool")
	if s.perfToolLn != nil {
		s.perfToolLn.Close()
		s.perfToolLn = nil
		s.dev.ReverseKillLocalAbstract(perfSockNamePre)
	}
	if s.frameSocket != nil {
		s.frameSocket.Close()
	}
}
