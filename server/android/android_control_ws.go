package android

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/touch"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AndroidControl(r *gin.Engine) {
	r.GET("/android/control", func(c *gin.Context) {

		var touchMap map[string]*touch.Touch = make(map[string]*touch.Touch)

		var perfMap map[string]context.CancelFunc = make(map[string]context.CancelFunc)

		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		serialInfo := &entity.SerialInfo{}
		// todo add error
		ws.ReadJSON(serialInfo)

		device, err := android_util.GetDevice(client, serialInfo.SerialName)
		if err != nil {
			ws.WriteJSON(entity.NewPerfDataError(err.Error()))
			log.Error(err)
		}

		if touchMap[device.Serial()] == nil {
			touchMap[device.Serial()] = touch.NewTouch(device)
		}

		control := touchMap[device.Serial()]

		exitCtx, exitFn := context.WithCancel(context.Background())

		var message entity.PerfRecvMessage

		go func() {
			for {
				select {
				case <-exitCtx.Done():
					delete(touchMap, device.Serial())
					return
				default:
					defer func() {
						if r := recover(); r != nil {
							log.Error("android control ws recovered:", r)
							exitFn()
						}
					}()
					err := ws.ReadJSON(&message)
					if err != nil {
						log.Error("android control read message steam err:", err)
						ws.WriteJSON(entity.NewPerfDataError("android control read message steam err:" + err.Error()))
						return
					} else {
						data, err1 := json.Marshal(message.Data)
						if err1 != nil {
							log.Error("control the data sent is not json")
							ws.WriteJSON(entity.NewPerfDataError("control the data sent is not json"))
							return
						}
						switch message.MessageType {
						case entity.ClosePerfType:
							log.Println("client send close perf info,close perf...")
							if perfExitFn, ok := perfMap[device.Serial()]; ok {
								perfExitFn()
							}
						case entity.StartPerfType:
							var perfConfig = &entity.PerfConfig{
								IntervalTime: 1,
							}

							perfExitCtx, perfExitFn := context.WithCancel(exitCtx)

							perfMap[device.Serial()] = perfExitFn

							err1 = json.Unmarshal(data, perfConfig)
							if err1 != nil {
								break
							}
							perfConfig.Ctx = perfExitCtx
							perfConfig.CancelFn = perfExitFn
							initPerfAndStart(serialInfo, perfConfig, device, ws)
						case entity.ControlTouchType:
							touchInfo := &entity.TouchInfo{}
							err1 = json.Unmarshal(data, touchInfo)
							if err1 != nil {
								break
							}
							control.Touch(*touchInfo)
						}

						if err1 != nil {
							log.Error("conversion message error,", err1)
							ws.WriteJSON(entity.NewPerfDataError(err1.Error()))
						}
					}
				}
			}
		}()
	})
}
