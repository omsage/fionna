package perf

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fionna/android/gadb"
	"fionna/entity"
)

func getMemInfo(device *gadb.Device, stats *entity.SystemInfo) (err error) {
	data, err := device.RunShellCommandWithBytes("cat /proc/meminfo")
	stats.MemInfo.Timestamp = time.Now().UnixMilli()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			switch parts[0] {
			case "MemTotal:":
				stats.MemInfo.MemTotal = val / 1024
			case "MemFree:":
				stats.MemInfo.MemFree = val / 1024
			case "Buffers:":
				stats.MemInfo.MemBuffers = val / 1024
			case "Cached:":
				stats.MemInfo.MemCached = val / 1024
			case "SwapTotal:":
				stats.MemInfo.SwapTotal = val / 1024
			case "SwapFree:":
				stats.MemInfo.SwapFree = val / 1024
			}
		}
	}
	stats.MemInfo.MemUsage = (stats.MemInfo.MemTotal - stats.MemInfo.MemFree - stats.MemInfo.MemBuffers - stats.MemInfo.MemCached) / 1024
	return
}

func getTemperatureInfo(device *gadb.Device, temperature *entity.SysTemperature) error {
	data, err := device.RunShellCommandWithBytes("dumpsys battery")
	if err != nil {
		return err
	}
	temperaturePattern := regexp.MustCompile(`temperature:\s*(\d+)`)
	match := temperaturePattern.FindStringSubmatch(string(data))

	if len(match) > 1 {
		temperatureStr := match[1]
		num, err := strconv.Atoi(strings.TrimSpace(temperatureStr))
		if err != nil {
			return err
		}
		temperature.Temperature = float64(num) / 10
	} else {
		return errors.New("temperature not found")
	}
	return nil
}

func getInterfaces(device *gadb.Device, stats *entity.SystemInfo) (err error) {
	data, err := device.RunShellCommandWithBytes("ip -o addr")
	if err != nil {
		// try /sbin/ip
		data, err = device.RunShellCommandWithBytes("/lib/ip -o addr")
		if err != nil {
			return
		}
	}

	if stats.NetworkInfo == nil {
		stats.NetworkInfo = make(map[string]*entity.SystemNetworkInfo)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := stats.NetworkInfo[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetworkInfo[intfname] = info
			} else {
				info := &entity.SystemNetworkInfo{
					InterfaceName: intfname,
				}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetworkInfo[intfname] = info
			}
		}
	}

	return
}

func getInterfaceInfo(device *gadb.Device, stats *entity.SystemInfo) (err error) {
	data, err := device.RunShellCommandWithBytes("cat /proc/net/dev")
	if err != nil {
		return
	}

	if stats.NetworkInfo == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := stats.NetworkInfo[intf]; ok {
				rx, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					continue
				}
				tx, err := strconv.ParseUint(parts[9], 10, 64)
				if err != nil {
					continue
				}
				info.Rx = float64(rx) / 1024 / 1024
				info.Tx = float64(tx) / 1024 / 1024
				info.Timestamp = time.Now().UnixMilli()
				stats.NetworkInfo[intf] = info
			}
		}
	}
	return
}

func parseCPUFields(fields []string, stat *entity.SystemCpuRaw) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}

type sysPreCpuInfo struct {
	// the CPU stats that were fetched last time round
	preCPU    entity.SystemCpuRaw
	preCPUMap map[string]entity.SystemCpuRaw
}

func getCPU(device *gadb.Device, stats *entity.SystemInfo, sysPreInfo *sysPreCpuInfo) (err error) {
	data, err := device.RunShellCommandWithBytes("cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU entity.SystemCpuRaw
		total  float32
	)

	//if preCPUMap == nil {
	//	preCPUMap = make(map[string]entity.SystemCpuRaw)
	//}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") { // changing here if you want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
			sysPreInfo.preCPU = sysPreInfo.preCPUMap[fields[0]]
			if sysPreInfo.preCPU.Total == 0 { // having no pre raw cpu data
				sysPreInfo.preCPUMap[fields[0]] = nowCPU
				continue
			}

			total = float32(nowCPU.Total - sysPreInfo.preCPU.Total)
			if stats.CPU == nil {
				stats.CPU = make(map[string]*entity.SystemCPUInfo)
			}
			cpu := &entity.SystemCPUInfo{}
			cpu.User = float32(nowCPU.User-sysPreInfo.preCPU.User) / total * 100
			cpu.Nice = float32(nowCPU.Nice-sysPreInfo.preCPU.Nice) / total * 100
			cpu.System = float32(nowCPU.System-sysPreInfo.preCPU.System) / total * 100
			cpu.Idle = float32(nowCPU.Idle-sysPreInfo.preCPU.Idle) / total * 100
			cpu.Iowait = float32(nowCPU.Iowait-sysPreInfo.preCPU.Iowait) / total * 100
			cpu.Irq = float32(nowCPU.Irq-sysPreInfo.preCPU.Irq) / total * 100
			cpu.SoftIrq = float32(nowCPU.SoftIrq-sysPreInfo.preCPU.SoftIrq) / total * 100
			cpu.Guest = float32(nowCPU.Guest-sysPreInfo.preCPU.Guest) / total * 100
			var cpuNowTime = float32(nowCPU.User + nowCPU.Nice + nowCPU.System + nowCPU.Iowait + nowCPU.Irq + nowCPU.SoftIrq)
			var cpuPreTime = float32(sysPreInfo.preCPU.User + sysPreInfo.preCPU.Nice + sysPreInfo.preCPU.System + sysPreInfo.preCPU.Iowait + sysPreInfo.preCPU.Irq + sysPreInfo.preCPU.SoftIrq)

			cpu.Usage = (cpuNowTime - cpuPreTime) / ((cpuNowTime + float32(nowCPU.Idle)) - (cpuPreTime + float32(sysPreInfo.preCPU.Idle))) * 100
			cpu.Timestamp = time.Now().UnixMilli()
			cpu.CPUName = fields[0]
			stats.CPU[fields[0]] = cpu
		}
	}
	return nil
}

func GetSysFrame(device *gadb.Device, perfOption entity.PerfConfig, frameCallBackFn func(frame *entity.SysFrameInfo, code entity.ServerCode)) {
	framePerf := NewPerfTool(device, perfOption.Ctx)

	framePerf.Init()
	framePerf.GetFrame(frameCallBackFn)
}

func GetSysCPU(device *gadb.Device, perfOption entity.PerfConfig, sysCpuCallBackFn func(CPU map[string]*entity.SystemCPUInfo, code entity.ServerCode)) {
	preCpuInfo := &sysPreCpuInfo{
		preCPUMap: make(map[string]entity.SystemCpuRaw),
		preCPU:    entity.SystemCpuRaw{},
	}
	time.Sleep(time.Duration(perfOption.IntervalTime * int(time.Second)))
	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	isNoFirst := false
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					systemInfo := &entity.SystemInfo{}
					err := getCPU(device, systemInfo, preCpuInfo)
					if err != nil {
						//systemInfo.Error = append(systemInfo.Error, err.Error())
						log.Error(err)
						sysCpuCallBackFn(nil, entity.GetPerfErr)
						return
					}
					if isNoFirst {
						sysCpuCallBackFn(systemInfo.CPU, entity.RequestSucceed)
					}
					isNoFirst = true
				}()
			}
		}
	}()
}

func GetSysMem(device *gadb.Device, perfOption entity.PerfConfig, sysMemCallBackFn func(sysMem *entity.SystemMemInfo, code entity.ServerCode)) {

	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					systemInfo := &entity.SystemInfo{}
					systemInfo.MemInfo = &entity.SystemMemInfo{}
					err := getMemInfo(device, systemInfo)
					if err != nil {
						//systemInfo.Error = append(systemInfo.Error, err.Error())
						log.Error(err)
						sysMemCallBackFn(nil, entity.GetPerfErr)
						return
					}
					sysMemCallBackFn(systemInfo.MemInfo, entity.RequestSucceed)
				}()

			}
		}
	}()
}

func GetSysNetwork(device *gadb.Device, perfOption entity.PerfConfig, sysNetworkCallBackFn func(sysNet map[string]*entity.SystemNetworkInfo, code entity.ServerCode)) {

	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					systemInfo := &entity.SystemInfo{}
					err := getInterfaces(device, systemInfo)
					if err != nil {
						//systemInfo.Error = append(systemInfo.Error, err.Error())
						log.Error(err)
						sysNetworkCallBackFn(nil, entity.GetPerfErr)
						return
					}
					err = getInterfaceInfo(device, systemInfo)
					if err != nil {
						//systemInfo.Error = append(systemInfo.Error, err.Error())
						log.Error(err)
						sysNetworkCallBackFn(nil, entity.GetPerfErr)
						return
					}
					sysNetworkCallBackFn(systemInfo.NetworkInfo, entity.RequestSucceed)
				}()

			}
		}
	}()
	return
}

func GetSysTemperature(device *gadb.Device, perfOption entity.PerfConfig, sysTemperatureCallBackFn func(sysTemperature *entity.SysTemperature, code entity.ServerCode)) {
	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					systemInfo := &entity.SysTemperature{
						Timestamp: time.Now().UnixMilli(),
					}
					err := getTemperatureInfo(device, systemInfo)
					if err != nil {
						log.Error(err)
						sysTemperatureCallBackFn(nil, entity.GetPerfErr)
						return
					}
					sysTemperatureCallBackFn(systemInfo, entity.RequestSucceed)
				}()

			}
		}
	}()
	return
}
