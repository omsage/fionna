package android

import (
	"context"
	"fionna/android/android_util"
	"fionna/android/touch"
	"fionna/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Android_Control(r *gin.Engine) {
	r.GET("/android/control", func(c *gin.Context) {
		serial := c.Query("udid")
		if serial == "" {
			log.Error("serial is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "serial is empty",
				Code: entity.ParameterErr,
			})
			return
		}
		device, err := android_util.GetDevice(client, serial)
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: entity.CodeDefaultMessage[entity.GetDeviceErr],
				Code: entity.GetDeviceErr,
			})
			log.Error(err)
		}

		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		control := touch.NewTouch(device)

		exitCtx, exitFn := context.WithCancel(context.Background())

		var message entity.TouchInfo

		go func() {
			for {
				select {
				case <-exitCtx.Done():
					return
				default:
					defer func() {
						if r := recover(); r != nil {
							log.Error("android control ws recovered:", r)
							exitFn()
						}
					}()
					err := ws.ReadJSON(&message)
					fmt.Println(message)
					if err != nil {
						log.Error("android control read message steam err:", err)
						ws.WriteJSON(entity.NewPerfDataError("android control read message steam err:" + err.Error()))
						break
					} else {
						control.Touch(message)
					}
				}
			}
		}()
	})
}
