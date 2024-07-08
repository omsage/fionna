package server

import (
	"fionna/entity"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
)

func GetReportInfoByPage(p *Pagination, name string) (reportList []entity.SerialInfo, err error) {
	err = db.Model(&entity.SerialInfo{}).Where("test_name like ?", "%"+name+"%").Order("created_at desc").Scopes(p.GormPaginate()).Find(&reportList).Error
	if err != nil {
		return nil, err
	}
	var total int64
	db.Model(&entity.SerialInfo{}).Count(&total)
	p.Total = cast.ToInt(total)
	return
}

func GroupReportUrl(r *gin.Engine) {
	reportUrlGroup := r.Group("/report")

	reportUrlGroup.GET("/list", func(c *gin.Context) {

		name := c.Query("name")

		p := NewPagination(c)

		reportList, err := GetReportInfoByPage(p, name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": "DB Error"})
			return
		}

		c.JSON(http.StatusOK, entity.ResponseData{
			Data: map[string]interface{}{
				"page":    p.Page,
				"size":    p.Size,
				"total":   p.Total,
				"reports": reportList,
			},
			Code: entity.RequestSucceed,
		})

	})

	reportUrlGroup.GET("/down", func(c *gin.Context) {

		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		serialInfo := &entity.SerialInfo{}

		db.First(serialInfo, "uuid = ?", uuid)

		c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
		//设置文件名称
		c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(*serialInfo.TestName)+".xlsx")
		file := Export2Excel(uuid)
		buffer, _ := file.WriteToBuffer()
		_, _ = c.Writer.Write(buffer.Bytes())
	})

	reportUrlGroup.POST("/rename", func(c *gin.Context) {

		info := entity.SerialInfo{}

		if err := c.ShouldBindJSON(&info); err != nil {
			// 如果解析失败，则返回错误响应
			log.Error("rename data is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "rename data is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		// User 的 ID 是 `111`
		db.Model(&info).Update("test_name", info.TestName)

		c.JSON(http.StatusOK, entity.ResponseData{
			Data: "rename succeed",
			Code: entity.RequestSucceed,
		})

	})

	reportUrlGroup.POST("/delete", func(c *gin.Context) {
		// 只要传一个uuid的json数组就行了
		infos := []string{}

		if err := c.ShouldBindJSON(&infos); err != nil {
			// 如果解析失败，则返回错误响应
			log.Error("delete report data is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "rename data is empty",
				Code: entity.ParameterErr,
			})
			return
		}
		// 可以组合成一数组的对应model删除，但是太麻烦。。。。
		for _, info := range infos {
			var perfConfig entity.PerfConfig
			db.First(&perfConfig, "uuid = ?", info)

			if perfConfig.FPS || perfConfig.Jank {
				db.Delete(&entity.SysFrameInfo{UUID: info})
				db.Delete(&entity.FrameSummary{UUID: info})
			}

			if perfConfig.SysCpu {
				db.Delete(&entity.SystemCPUInfo{UUID: info})
				db.Delete(&entity.SystemCPUSummary{UUID: info})
			}

			if perfConfig.SysMem {
				db.Delete(&entity.SystemMemInfo{UUID: info})
				db.Delete(&entity.SystemMemSummary{UUID: info})
			}

			if perfConfig.SysTemperature {
				db.Delete(&entity.SystemTemperatureSummary{UUID: info})
				db.Delete(&entity.SysTemperature{UUID: info})
			}

			if perfConfig.SysNetwork {
				db.Delete(&entity.SystemNetworkInfo{UUID: info})
				db.Delete(&entity.SystemNetworkSummary{UUID: info})
			}

			if perfConfig.ProcCpu {
				db.Delete(&entity.ProcCpuInfo{UUID: info})
				db.Delete(&entity.ProcCpuSummary{UUID: info})
			}

			if perfConfig.ProcMem {
				db.Delete(&entity.ProcMemInfo{UUID: info})
				db.Delete(&entity.ProcMemSummary{UUID: info})
			}

			if perfConfig.ProcThread {
				db.Delete(&entity.ProcThreadsInfo{UUID: info})
			}

			db.Delete(perfConfig)

			db.Delete(&entity.SerialInfo{UUID: info})
		}

	})

	reportUrlGroup.GET("/config", func(c *gin.Context) {

		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		c.JSON(http.StatusOK, entity.ResponseData{
			Data: perfConfig,
			Code: entity.RequestSucceed,
		})

	})

	reportUrlGroup.GET("/summary", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		overallSummary := &entity.OverallSummary{}

		if perfConfig.FPS || perfConfig.Jank {
			var frameSummary entity.FrameSummary
			db.First(&frameSummary, "uuid = ?", uuid)
			overallSummary.SysFrameSummary = &frameSummary
		}

		if perfConfig.SysCpu {
			sysCpuSummary := make(map[string]entity.SystemCPUSummary)
			var cpuSummarys []entity.SystemCPUSummary
			db.Where("uuid = ?", uuid).Find(&cpuSummarys)
			for _, value := range cpuSummarys {
				sysCpuSummary[value.CpuName] = value
			}
			overallSummary.SysCpuSummary = sysCpuSummary
		}

		if perfConfig.SysNetwork {
			sysNetworkSummary := make(map[string]entity.SystemNetworkSummary)
			var netSummarys []entity.SystemNetworkSummary
			db.Where("uuid = ?", uuid).Find(&netSummarys)
			for _, value := range netSummarys {
				sysNetworkSummary[value.Name] = value
			}
			overallSummary.NetworkSummary = sysNetworkSummary
		}

		if perfConfig.SysMem {
			var sysMemSummarys entity.SystemMemSummary
			db.First(&sysMemSummarys, "uuid = ?", uuid)
			overallSummary.SysMemSummary = &sysMemSummarys
		}

		if perfConfig.SysTemperature {
			var sysTemperature entity.SystemTemperatureSummary
			db.First(&sysTemperature, "uuid = ?", uuid)
			overallSummary.SysTemperatureSummary = &sysTemperature
		}

		if perfConfig.ProcCpu {
			var procCpuSummary entity.ProcCpuSummary
			db.First(&procCpuSummary, "uuid = ?", uuid)
			overallSummary.ProcCpu = &procCpuSummary
		}

		if perfConfig.ProcMem {
			var procMemSummary entity.ProcMemSummary
			db.First(&procMemSummary, "uuid = ?", uuid)
			overallSummary.ProcMem = &procMemSummary
		}

		c.JSON(http.StatusOK, entity.ResponseData{
			Data: overallSummary,
			Code: entity.RequestSucceed,
		})

	})

	reportUrlGroup.GET("/proc/cpu", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.ProcCpu {
			var procCpuDatas []entity.ProcCpuInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procCpuDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: procCpuDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "process performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/proc/mem", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.ProcMem {
			var procMemDatas []entity.ProcMemInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procMemDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: procMemDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "process cpu mem performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/proc/thread", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.ProcThread {
			var procThreadDatas []entity.ProcThreadsInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procThreadDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: procThreadDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "process thread performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/sys/cpu", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.SysCpu {
			var sysCpuDatas []entity.SystemCPUInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysCpuDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: sysCpuDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "sys cpu usage performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/sys/mem", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.SysMem {
			var SYSMemDatas []entity.SystemMemInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&SYSMemDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: SYSMemDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "sys mem performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/sys/frame", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.FPS || perfConfig.Jank {
			var sysFrameDatas []entity.SysFrameInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysFrameDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: sysFrameDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "sys frame performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/sys/network", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.SysNetwork {
			var SysNetworkDatas []entity.SystemNetworkInfo
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&SysNetworkDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: SysNetworkDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "system network performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

	reportUrlGroup.GET("/sys/temperature", func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			log.Error("uuid is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "uuid is empty",
				Code: entity.ParameterErr,
			})
			return
		}

		var perfConfig entity.PerfConfig
		db.First(&perfConfig, "uuid = ?", uuid)

		if perfConfig.SysNetwork {
			var sysTemperatureDatas []entity.SysTemperature
			db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysTemperatureDatas)
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: sysTemperatureDatas,
				Code: entity.RequestSucceed,
			})
		} else {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "system temperature performance was not collected",
				Code: entity.GetPerfErr,
			})
		}

	})

}
