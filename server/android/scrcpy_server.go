package android

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/scrcpy_client"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func WebSocketScrcpy(r *gin.Engine) {
	r.GET("/android/scrcpy", func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}
		//defer ws.Close()

		var udidInfo entity.ScrcpyDevice

		err = ws.ReadJSON(&udidInfo)

		if err != nil {
			log.Error("get scrcpy id error,", err)
			ws.WriteJSON(entity.NewScrcpyError("get scrcpy id error," + err.Error()))
		}
		dev, err := android_util.GetDevice(client, udidInfo.UDID)
		if err != nil {
			ws.WriteJSON(entity.NewScrcpyError(err.Error()))
			log.Error(err)
			return
		}
		exitCtx, _ := context.WithCancel(context.Background())

		scrcpyClient := scrcpy_client.NewScrcpy(dev, exitCtx, ws)

		scrcpyClient.Start()

		var message entity.ScrcpyRecvMessage

		go func() {
			for {
				select {
				case <-exitCtx.Done():
					return
				default:
					defer func() {
						if r := recover(); r != nil {
							log.Error("ws recovered:", r)
							scrcpyClient.ClientStop()
						}
					}()
					err := ws.ReadJSON(&message)
					if err != nil {
						log.Error("read message steam err:", err)
						ws.WriteJSON(entity.NewScrcpyError("read message steam err:" + err.Error()))
						scrcpyClient.ClientStop()
					} else {
						if message.MessageType == entity.ScrcpyTouchType {
							data, err1 := json.Marshal(message.Data)
							if err1 != nil {
								log.Error("the data sent is not json")
								ws.WriteJSON(entity.NewScrcpyError("the data sent is not json"))
								break
							}
							var touch = &entity.ScrcpyTouch{}
							err1 = json.Unmarshal(data, touch)
							if err1 == nil {
								err = scrcpyClient.Touch(touch, 1)
								if err != nil {
									log.Error("execute touch err:", err)
									ws.WriteJSON(entity.NewScrcpyError("execute touch err:" + err.Error()))
									break
								}
							} else {
								log.Error("conversion message error,", err1)
								ws.WriteJSON(entity.NewScrcpyError("conversion message error," + err1.Error()))
								break
							}
						}
						if message.MessageType == entity.ScrcpyCloseType {
							scrcpyClient.ClientStop()
						}
						if message.MessageType == entity.ScrcpyPongType {
							continue
						}
					}
				}
			}
		}()

	})

}
