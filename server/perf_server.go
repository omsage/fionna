package server

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

func frameDataSummary(frameOverview *entity.FrameSummary, frame *entity.SysFrameInfo, count float64) {
	frameOverview.AvgFPS = (frameOverview.AvgFPS*count + float64(frame.FPS)) / (count + 1)
	frameOverview.AllJankCount = frameOverview.AllJankCount + frame.JankCount
	frameOverview.AllBigJankCount = frameOverview.AllBigJankCount + frame.BigJankCount
	frameOverview.JankCountRate = float64(frameOverview.AllJankCount) / (count + 1) * 100
	frameOverview.BigJankCountRate = float64(frameOverview.AllBigJankCount) / (count + 1) * 100
	if frameOverview.MaxBigJankCount < frame.BigJankCount {
		frameOverview.MaxBigJankCount = frame.BigJankCount
	}
	if frameOverview.MaxJankCount < frame.JankCount {
		frameOverview.MaxJankCount = frame.JankCount
	}
}

func sysCpuDataSummary(sysCpuOverview *entity.SystemCPUSummary, sysCPU *entity.SystemCPUInfo, count float64) {
	sysCpuOverview.AvgSysCPU = (sysCpuOverview.AvgSysCPU*count + float64(sysCPU.Usage)) / (count + 1)
	if sysCpuOverview.MaxSysCPU < float64(sysCPU.Usage) {
		sysCpuOverview.MaxSysCPU = float64(sysCPU.Usage)
	}
}

func sysMemDataSummary(sysMemOverview *entity.SystemMemSummary, sysMem *entity.SystemMemInfo, count float64) {
	sysMemOverview.AvgMemTotal = (sysMemOverview.AvgMemTotal*count + float64(sysMem.MemTotal)) / (count + 1)
	if sysMemOverview.MaxMemTotal < float64(sysMem.MemTotal) {
		sysMemOverview.MaxMemTotal = float64(sysMem.MemTotal)
	}
}

func sysTemperatureDataSummary(sysTemperatureOverview *entity.SystemTemperatureSummary, sysTemperature *entity.SysTemperature) {
	if sysTemperatureOverview.MaxTemperature < sysTemperature.Temperature {
		sysTemperatureOverview.MaxTemperature = sysTemperature.Temperature
	}
}

func procCpuDataSummary(procCpuOverview *entity.ProcCpuSummary, procCpu *entity.ProcCpuInfo, count float64) {
	procCpuOverview.AvgProcCPU = (procCpuOverview.AvgProcCPU*count + procCpu.CpuUtilization) / (count + 1)
	if procCpuOverview.MaxProcCPU < procCpu.CpuUtilization {
		procCpuOverview.MaxProcCPU = procCpu.CpuUtilization
	}
}

func procMemDataSummary(procMemOverview *entity.ProcMemSummary, procMem *entity.ProcMemInfo, count float64) {
	procMemOverview.AvgTotalPSS = (procMemOverview.AvgTotalPSS*count + procMem.TotalPSS) / (count + 1)
	procMemOverview.AvgCode = (procMemOverview.AvgCode*count + procMem.Code) / (count + 1)
	procMemOverview.AvgGraphics = (procMemOverview.AvgGraphics*count + procMem.Graphics) / (count + 1)
	procMemOverview.AvgJavaHeap = (procMemOverview.AvgJavaHeap*count + procMem.JavaHeap) / (count + 1)
	procMemOverview.AvgNativeHeap = (procMemOverview.AvgNativeHeap*count + procMem.NativeHeap) / (count + 1)
	procMemOverview.AvgPrivateOther = (procMemOverview.AvgPrivateOther*count + procMem.PrivateOther) / (count + 1)
	procMemOverview.AvgStack = (procMemOverview.AvgStack*count + procMem.Stack) / (count + 1)
	procMemOverview.AvgSystem = (procMemOverview.AvgSystem*count + procMem.System) / (count + 1)

	if procMemOverview.MaxTotalPSS < procMem.TotalPSS {
		procMemOverview.MaxTotalPSS = procMem.TotalPSS
	}
	if procMemOverview.MaxCode < procMem.Code {
		procMemOverview.MaxCode = procMem.Code
	}
	if procMemOverview.MaxGraphics < procMem.Graphics {
		procMemOverview.MaxGraphics = procMem.Graphics
	}
	if procMemOverview.MaxJavaHeap < procMem.JavaHeap {
		procMemOverview.MaxJavaHeap = procMem.JavaHeap
	}
	if procMemOverview.MaxNativeHeap < procMem.NativeHeap {
		procMemOverview.MaxNativeHeap = procMem.NativeHeap
	}
	if procMemOverview.MaxPrivateOther < procMem.PrivateOther {
		procMemOverview.MaxPrivateOther = procMem.PrivateOther
	}
	if procMemOverview.MaxStack < procMem.Stack {
		procMemOverview.MaxStack = procMem.Stack
	}
	if procMemOverview.MaxSystem < procMem.System {
		procMemOverview.MaxSystem = procMem.System
	}
}

func startGetPerf(perfWsConn *websocket.Conn, device *gadb.Device, config entity.PerfConfig) {
	perfConn := NewSafeWebsocket(perfWsConn)
	if config.FPS || config.Jank {

		frameOverview := entity.NewFrameSummary(config.UUID)

		count := 0.0

		go func() {
			perf.GetSysFrame(device, config, func(frame *entity.SysFrameInfo, code entity.ServerCode) {

				sysFrameInfo := &entity.SysFrameInfo{
					UUID:      config.UUID,
					Timestamp: frame.Timestamp,
				}

				if config.FPS {
					sysFrameInfo.FPS = frame.FPS
				}
				if config.Jank {
					sysFrameInfo.JankCount = frame.JankCount
					sysFrameInfo.BigJankCount = frame.BigJankCount
				}

				frameDataSummary(frameOverview, sysFrameInfo, count)

				count++
				go func() {
					if count == 1 {
						db.Create(frameOverview)
					} else {
						db.Save(frameOverview)
					}

					db.Create(sysFrameInfo)
				}()

				perfData := &entity.PerfData{SystemPerfData: &entity.SystemInfo{Frame: sysFrameInfo}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send sys frame fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}
	if config.SysCpu {

		systemCPUOverviewInfo := make(map[string]*entity.SystemCPUSummary)

		count := 0.0

		go func() {
			perf.GetSysCPU(device, config, func(CPU map[string]*entity.SystemCPUInfo, code entity.ServerCode) {

				for cpuName, value := range CPU {
					value.UUID = config.UUID

					systemCPUOverview := systemCPUOverviewInfo[value.CPUName]

					if systemCPUOverview == nil {
						systemCPUOverview = entity.NewSystemCPUSummary(config.UUID)
						systemCPUOverview.CpuName = cpuName
						systemCPUOverviewInfo[value.CPUName] = systemCPUOverview
					}
					sysCpuDataSummary(systemCPUOverview, value, count)

					go func() {
						db.Create(value)
						if count == 0 {
							db.Create(systemCPUOverview)
						} else {
							db.Save(systemCPUOverview)
						}
					}()

					count++
				}

				perfData := &entity.PerfData{SystemPerfData: &entity.SystemInfo{CPU: CPU}}

				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send sys cpu fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}
	if config.SysMem {
		count := 0.0
		sysMemOverview := entity.NewSystemMemSummary(config.UUID)
		go func() {
			perf.GetSysMem(device, config, func(sysMem *entity.SystemMemInfo, code entity.ServerCode) {
				sysMem.UUID = config.UUID

				sysMemDataSummary(sysMemOverview, sysMem, count)

				go func() {
					db.Create(sysMem)
					if count == 0 {
						db.Create(sysMemOverview)
					} else {
						db.Save(sysMemOverview)
					}
				}()

				count++

				perfData := &entity.PerfData{SystemPerfData: &entity.SystemInfo{MemInfo: sysMem}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send sys mem fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}
	if config.SysNetwork {

		count := 0.0

		sysNetInit := make(map[string]*entity.SystemNetworkInfo)
		sysNetOverviews := make(map[string]*entity.SystemNetworkSummary)

		go func() {
			perf.GetSysNetwork(device, config, func(sysNet map[string]*entity.SystemNetworkInfo, code entity.ServerCode) {
				for name, value := range sysNet {

					if count == 0 {
						sysNetInit[value.InterfaceName] = &entity.SystemNetworkInfo{
							Rx: value.Rx,
							Tx: value.Tx,
						}
						sysNetOverviews[value.InterfaceName] = entity.NewSystemNetworkSummary(config.UUID, name)
					}

					initNetwork := sysNetInit[value.InterfaceName]
					value.UUID = config.UUID
					value.Rx = value.Rx - initNetwork.Rx
					value.Tx = value.Tx - initNetwork.Tx

					sysNetOverview := sysNetOverviews[name]

					//if strings.Contains(value.InterfaceName, "wlan") {

					sysNetOverview.AllSysTxData += value.Tx
					sysNetOverview.AllSysRxData += value.Rx

					go func() {
						db.Create(value)
						if count == 0 {
							db.Create(sysNetOverview)
						} else {
							db.Save(sysNetOverview)
						}
					}()

				}
				count++
				perfData := &entity.PerfData{SystemPerfData: &entity.SystemInfo{NetworkInfo: sysNet}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send sys network fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}

	if config.ProcCpu {

		count := 0.0
		procCpuOverview := entity.NewProcCpuSummary(config.UUID)

		go func() {
			perf.GetProcCPU(device, config, func(cpuInfo *entity.ProcCpuInfo, code entity.ServerCode) {
				cpuInfo.UUID = config.UUID

				procCpuDataSummary(procCpuOverview, cpuInfo, count)

				go func() {
					db.Create(cpuInfo)
					if count == 0 {
						db.Create(procCpuOverview)
					} else {
						db.Save(procCpuOverview)
					}
				}()

				count++

				perfData := &entity.PerfData{ProcPerfData: &entity.ProcessInfo{CPUInfo: cpuInfo}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send proc cpu fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}

	if config.ProcMem {
		count := 0.0
		procMemOverview := entity.NewProcMemSummary(config.UUID)
		go func() {
			perf.GetProcMem(device, config, func(memInfo *entity.ProcMemInfo, code entity.ServerCode) {

				memInfo.UUID = config.UUID

				procMemDataSummary(procMemOverview, memInfo, count)
				//
				//data, _ := json.Marshal(procMemOverview)
				//
				//fmt.Println(string(data))

				go func() {
					db.Create(memInfo)
					if count == 0 {
						db.Create(procMemOverview)
					} else {
						db.Save(procMemOverview)
					}
				}()

				count++

				perfData := &entity.PerfData{ProcPerfData: &entity.ProcessInfo{MemInfo: memInfo}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send proc mem fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()

	}

	if config.ProcThread {
		go func() {
			perf.GetProcThreads(device, config, func(threadInfo *entity.ProcThreadsInfo, code entity.ServerCode) {

				threadInfo.UUID = config.UUID
				go func() {
					db.Create(threadInfo)
				}()

				perfData := &entity.PerfData{ProcPerfData: &entity.ProcessInfo{ThreadInfo: threadInfo}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send proc thread fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()
	}

	if config.SysTemperature {

		sysTemperatureSummary := entity.NewSystemTemperature(config.UUID)
		count := 0.0

		var initTemperature = 0.0

		go func() {
			perf.GetSysTemperature(device, config, func(temperatureInfo *entity.SysTemperature, code entity.ServerCode) {

				temperatureInfo.UUID = config.UUID

				sysTemperatureDataSummary(sysTemperatureSummary, temperatureInfo)

				if count == 0 {
					initTemperature = temperatureInfo.Temperature
				}

				sysTemperatureSummary.DiffTemperature = sysTemperatureSummary.MaxTemperature - initTemperature

				go func() {
					db.Create(temperatureInfo)
					if count == 0 {
						db.Create(sysTemperatureSummary)
					} else {
						db.Save(sysTemperatureSummary)
					}
				}()

				count++

				perfData := &entity.PerfData{SystemPerfData: &entity.SystemInfo{Temperature: temperatureInfo}}
				err := perfConn.WriteJSON(entity.NewPerfDataMessage(perfData))
				if err != nil {
					log.Error("perf conn send sys temperature fail,close perf....", err)
					config.CancelFn()
				}
			})
		}()
	}
}

func WebSocketPerf(r *gin.Engine) {
	r.GET("/android/perf", func(c *gin.Context) {

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

		exitCtx, exitFn := context.WithCancel(context.Background())

		var message entity.PerfRecvMessage

		go func() {
			for {
				select {
				case <-exitCtx.Done():
					return
				default:
					defer func() {
						if r := recover(); r != nil {
							log.Error("perf ws recovered:", r)
							exitFn()
						}
					}()
					err := ws.ReadJSON(&message)
					if err != nil {
						log.Error("perf read message steam err:", err)
						ws.WriteJSON(entity.NewPerfDataError("perf read message steam err:" + err.Error()))
						break
					} else {
						if message.MessageType == entity.StartPerfType {

							data, err1 := json.Marshal(message.Data)
							if err1 != nil {
								log.Error("perf the data sent is not json")
								ws.WriteJSON(entity.NewPerfDataError("perf the data sent is not json"))
								break
							}
							// todo uuid
							var perfConfig = &entity.PerfConfig{
								IntervalTime: 1,
							}
							err1 = json.Unmarshal(data, perfConfig)

							if err1 == nil {

								id := uuid.New()

								//reportBase := &entity.BaseModel{
								//	UUID: id.String(),
								//}
								currentTime := time.Now()

								// 格式化时间为字符串
								formattedTime := currentTime.Format("2006-01-02 15:04:05")

								var testName = ""

								if perfConfig.PackageName != "" && perfConfig.Pid != "" {
									testName = fmt.Sprintf("%s_%s_%s_pid%s_%s", serialInfo.ProductDevice, serialInfo.Model, perfConfig.PackageName, perfConfig.Pid, formattedTime)
								} else if perfConfig.PackageName != "" && perfConfig.Pid == "" {

									testName = fmt.Sprintf("%s_%s_%s_%s", serialInfo.ProductDevice, serialInfo.Model, perfConfig.PackageName, formattedTime)

									perfConfig.Pid, err = android_util.GetPidOnPackageName(device, perfConfig.PackageName)

									if err != nil {
										log.Error("get pid err:", err)
										ws.WriteJSON(entity.NewPerfDataError("get pid err:" + err.Error()))
										break
									}
								} else if perfConfig.PackageName == "" && perfConfig.Pid != "" {
									testName = fmt.Sprintf("%s_%s_pid%s_%s", serialInfo.ProductDevice, serialInfo.Model, perfConfig.Pid, formattedTime)
								} else {
									testName = fmt.Sprintf("%s_%s_%s", serialInfo.ProductDevice, serialInfo.Model, formattedTime)
								}

								serialInfo.TestName = &testName
								timestamp := time.Now().UnixMilli()

								serialInfo.Timestamp = &timestamp
								serialInfo.PackageName = &perfConfig.PackageName

								if serialInfo.UUID == "" {
									serialInfo.UUID = id.String()
								}

								db.Create(serialInfo)

								if perfConfig.IntervalTime == 0 {
									perfConfig.IntervalTime = 1
								}

								perfConfig.Ctx = exitCtx
								perfConfig.CancelFn = exitFn

								perfConfig.UUID = id.String()
								db.Create(perfConfig)

								startGetPerf(ws, device, *perfConfig)

							} else {
								log.Error("conversion message error,", err1)
								ws.WriteJSON(entity.NewPerfDataError(err1.Error()))
								break
							}
						}
						if message.MessageType == entity.ClosePerfType {
							log.Println("client send close perf info,close perf...")
							exitFn()
						}
						if message.MessageType == entity.PongPerfType {
							continue
						}
					}
				}
			}
		}()
	})
}
