package android

import (
	"context"
	"fionna/android/android_util"
	"fionna/android/scrcpy_client"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WebSocketScrcpy(r *gin.Engine) {
	r.GET("/android/scrcpy", func(c *gin.Context) {

		pic := c.Query("pic")
		if pic == "" {
			log.Error("pic is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "pic is empty",
				Code: entity.ParameterErr,
			})
			return
		}

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

		scrcpyClient.Start(entity.ScrcpyPic(pic))

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
						return
					} else {
						if message.MessageType == entity.ScrcpyCloseType {
							scrcpyClient.ClientStop()
						}
						if message.MessageType == entity.ScrcpyPongType {
							continue
						}
						//if message.MessageType == entity.ScrcpyPicType {
						//	if pic, ok := message.Data.(entity.ScrcpyPic); ok {
						//		scrcpyClient.ClientStop()
						//		scrcpyClient.Start(pic)
						//	}
						//}
					}
				}
			}
		}()

	})

}
