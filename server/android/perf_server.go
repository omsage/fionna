package android

import (
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/android/perf"
	"fionna/entity"
	"fionna/server/db"
	"fionna/server/util"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
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
	perfConn := util.NewSafeWebsocket(perfWsConn)
	if config.FPS || config.Jank {

		frameOverview := entity.NewFrameSummary(config.UUID)

		count := 0.0

		go func() {
			var lock sync.Mutex
			perf.GetSysFrame(device, config, func(frame *entity.SysFrameInfo, code entity.ServerCode) {
				lock.Lock()
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

				go func() {

					if count == 1 {
						db.GetDB().Create(frameOverview)
					} else {
						db.GetDB().Save(frameOverview)
					}

					db.GetDB().Create(sysFrameInfo)

				}()
				count++
				lock.Unlock()
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
			var lock sync.Mutex
			perf.GetSysCPU(device, config, func(CPU map[string]*entity.SystemCPUInfo, code entity.ServerCode) {
				lock.Lock()
				go func() {
					d, _ := json.Marshal(CPU)
					sCpu := &entity.SystemCPUData{UUID: config.UUID, Data: string(d), Timestamp: time.Now().UnixMilli()}
					db.GetDB().Create(sCpu)
				}()
				for cpuName, value := range CPU {
					value.UUID = config.UUID

					systemCPUOverview := systemCPUOverviewInfo[value.CPUName]

					if systemCPUOverview == nil {
						systemCPUOverview = entity.NewSystemCPUSummary(config.UUID)
						systemCPUOverview.CpuName = cpuName
						systemCPUOverviewInfo[value.CPUName] = systemCPUOverview
					}
					sysCpuDataSummary(systemCPUOverview, value, count)

					go func(v entity.SystemCPUInfo) {

						db.GetDB().Create(&v)
						if count == 0 {
							db.GetDB().Create(systemCPUOverview)
						} else {
							db.GetDB().Save(systemCPUOverview)
						}

					}(*value)
				}
				count++
				lock.Unlock()
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
			var lock sync.Mutex
			perf.GetSysMem(device, config, func(sysMem *entity.SystemMemInfo, code entity.ServerCode) {
				sysMem.UUID = config.UUID
				lock.Lock()
				sysMemDataSummary(sysMemOverview, sysMem, count)

				go func() {

					db.GetDB().Create(sysMem)
					if count == 0 {
						db.GetDB().Create(sysMemOverview)
					} else {
						db.GetDB().Save(sysMemOverview)
					}

				}()

				count++
				lock.Unlock()
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
			var lock sync.Mutex
			perf.GetSysNetwork(device, config, func(sysNet map[string]*entity.SystemNetworkInfo, code entity.ServerCode) {
				lock.Lock()
				go func() {
					d, _ := json.Marshal(sysNet)
					sNet := &entity.SystemNetworkData{UUID: config.UUID, Data: string(d), Timestamp: time.Now().UnixMilli()}
					db.GetDB().Create(sNet)
				}()
				for name, netV := range sysNet {

					if count == 0 {
						sysNetInit[netV.InterfaceName] = &entity.SystemNetworkInfo{
							Rx: netV.Rx,
							Tx: netV.Tx,
						}
						sysNetOverviews[netV.InterfaceName] = entity.NewSystemNetworkSummary(config.UUID, name)
					}

					initNetwork := sysNetInit[netV.InterfaceName]
					netV.UUID = config.UUID
					netV.Rx = netV.Rx - initNetwork.Rx
					netV.Tx = netV.Tx - initNetwork.Tx

					sysNetOverview := sysNetOverviews[name]

					//if strings.Contains(netV.InterfaceName, "wlan") {

					sysNetOverview.AllSysTxData = netV.Tx
					sysNetOverview.AllSysRxData = netV.Rx

					go func() {

						if count == 0 {
							db.GetDB().Create(sysNetOverview)
						} else {
							db.GetDB().Save(sysNetOverview)
						}

					}()

				}
				count++
				lock.Unlock()
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
			var lock sync.Mutex
			perf.GetProcCPU(device, config, func(cpuInfo *entity.ProcCpuInfo, code entity.ServerCode) {
				cpuInfo.UUID = config.UUID
				lock.Lock()
				procCpuDataSummary(procCpuOverview, cpuInfo, count)

				go func() {

					db.GetDB().Create(cpuInfo)
					if count == 0 {
						db.GetDB().Create(procCpuOverview)
					} else {
						db.GetDB().Save(procCpuOverview)
					}

				}()

				count++
				lock.Unlock()
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
			var lock sync.Mutex
			perf.GetProcMem(device, config, func(memInfo *entity.ProcMemInfo, code entity.ServerCode) {

				memInfo.UUID = config.UUID
				lock.Lock()
				procMemDataSummary(procMemOverview, memInfo, count)
				//
				//data, _ := json.Marshal(procMemOverview)
				//
				//fmt.Println(string(data))

				go func() {

					db.GetDB().Create(memInfo)
					if count == 0 {
						db.GetDB().Create(procMemOverview)
					} else {
						db.GetDB().Save(procMemOverview)
					}

				}()

				count++
				lock.Unlock()
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
					db.GetDB().Create(threadInfo)
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
			var lock sync.Mutex
			perf.GetSysTemperature(device, config, func(temperatureInfo *entity.SysTemperature, code entity.ServerCode) {

				temperatureInfo.UUID = config.UUID

				sysTemperatureDataSummary(sysTemperatureSummary, temperatureInfo)
				lock.Lock()
				if count == 0 {
					initTemperature = temperatureInfo.Temperature
				}

				sysTemperatureSummary.DiffTemperature = sysTemperatureSummary.MaxTemperature - initTemperature

				go func() {

					db.GetDB().Create(temperatureInfo)
					if count == 0 {
						db.GetDB().Create(sysTemperatureSummary)
					} else {
						db.GetDB().Save(sysTemperatureSummary)
					}

				}()

				count++
				lock.Unlock()
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

func InitPerfAndStart(serialInfo *entity.SerialInfo, perfConfig *entity.PerfConfig, device *gadb.Device, ws *websocket.Conn) {
	id := uuid.New()

	currentTime := time.Now()

	// 格式化时间为字符串
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	var testName = ""
	var err error

	if perfConfig.PackageName != "" && perfConfig.Pid != "" {
		testName = fmt.Sprintf("%s_%s_%s_pid%s_%s", serialInfo.ProductDevice, serialInfo.Model, perfConfig.PackageName, perfConfig.Pid, formattedTime)
	} else if perfConfig.PackageName != "" && perfConfig.Pid == "" {

		testName = fmt.Sprintf("%s_%s_%s_%s", serialInfo.ProductDevice, serialInfo.Model, perfConfig.PackageName, formattedTime)

		perfConfig.Pid, err = android_util.GetPidOnPackageName(device, perfConfig.PackageName)

		if err != nil {
			log.Error("get pid err:", err)
			ws.WriteJSON(entity.NewPerfDataError("get pid err:" + err.Error()))
			return
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

	serialInfo.UUID = id.String()
	perfConfig.UUID = id.String()
	db.GetDB().Create(serialInfo)

	if perfConfig.IntervalTime == 0 {
		perfConfig.IntervalTime = 1
	}

	db.GetDB().Create(perfConfig)

	startGetPerf(ws, device, *perfConfig)
}
