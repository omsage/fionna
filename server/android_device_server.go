package server

import (
	"fionna/android/android_util"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GroupAndroidPackageUrl(r *gin.Engine) {
	packageUrlGroup := r.Group("/android/app")

	packageUrlGroup.GET("/list", func(c *gin.Context) {
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
		packageList, err := android_util.GetPackageNameList(device)
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "package list get error",
				Code: entity.CustomErr,
			})
			log.Error(err)
			return
		}
		c.JSON(http.StatusOK, entity.ResponseData{
			Data: packageList,
			Code: entity.RequestSucceed,
		})

	})

	packageUrlGroup.GET("/current", func(c *gin.Context) {
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
			if err != nil {
				c.JSON(http.StatusOK, entity.ResponseData{
					Data: entity.CodeDefaultMessage[entity.GetDeviceErr],
					Code: entity.GetDeviceErr,
				})
				log.Error(err)
				return
			}
			return
		}
		packageName, pid, err := android_util.GetCurrentPackageNameAndPid(device)
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "current package get error",
				Code: entity.CustomErr,
			})
			log.Error(err)
			return
		}
		c.JSON(http.StatusOK, entity.ResponseData{
			Data: map[string]string{
				"packageName": packageName,
				"pid":         pid,
			},
			Code: entity.RequestSucceed,
		})
	})

}

func GroupAndroidSerialUrl(r *gin.Engine) {

	serialUrlGroup := r.Group("/android/serial")

	serialUrlGroup.GET("/list", func(c *gin.Context) {
		serialList, err := android_util.GetSerialList(client)
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: entity.CodeDefaultMessage[entity.GetDeviceErr],
				Code: entity.GetDeviceErr,
			})
			log.Error(err)
			return
		}
		c.JSON(http.StatusOK, entity.ResponseData{
			Data: serialList,
			Code: entity.RequestSucceed,
		})
	})

	serialUrlGroup.GET("/default", func(c *gin.Context) {
		device, err := android_util.GetDevice(client, "")
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: entity.CodeDefaultMessage[entity.GetDeviceErr],
				Code: entity.GetDeviceErr,
			})
			log.Error(err)
		}
		c.JSON(http.StatusOK, entity.ResponseData{
			Data: device.Serial(),
			Code: entity.RequestSucceed,
		})
	})

	serialUrlGroup.GET("/info", func(c *gin.Context) {
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
		serialInfo := android_util.GetSerialInfo(device)
		c.JSON(http.StatusOK, entity.ResponseData{
			Data: serialInfo,
			Code: entity.RequestSucceed,
		})
	})

	serialUrlGroup.GET("/keycode", func(c *gin.Context) {
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

		keycode := c.Query("keycode")
		if keycode == "" {
			log.Error("keycode is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "keycode is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		device.RunShellCommand("input keyevent " + keycode)

		c.JSON(http.StatusOK, entity.ResponseData{
			Data: "",
			Code: entity.RequestSucceed,
		})
	})
}
