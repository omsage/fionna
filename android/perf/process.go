package perf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fionna/android/gadb"
	"fionna/entity"
)

type procCpuConfig struct {
	IntervalTime    int
	preProcCpuTime  float64
	preTotalCpuTime float64
}

func newProcCPUConfig() *procCpuConfig {
	return &procCpuConfig{preProcCpuTime: -1, preTotalCpuTime: -1}
}

func getStatOnPid(device *gadb.Device, pid string) (stat *entity.ProcessStat, err error) {
	data, err := device.RunShellCommand(fmt.Sprintf("cat /proc/%s/stat", pid))
	if err != nil {
		return nil, fmt.Errorf("exec command erro : " + fmt.Sprintf("cat /proc/%s/stat", pid))
	}
	return newProcessStat(data)
}

func getMemTotalPSS(device *gadb.Device, pid string) (memInfo *entity.ProcMemInfo, err error) {
	data, err := device.RunShellCommand(fmt.Sprintf("dumpsys meminfo --local %s", pid))
	if err != nil {
		return
	}
	if strings.Contains(data, "No process found for") {
		return nil, errors.New(data)
	}
	s := strings.Split(data, "\n")

	memInfo = &entity.ProcMemInfo{
		Timestamp: time.Now().UnixMilli(),
	}

	isGetMem := false

	for _, v := range s {
		v = strings.ReplaceAll(v, "\n", "")
		if strings.Contains(v, "App Summary") {
			isGetMem = true
			continue
		}

		if isGetMem && strings.Contains(v, ":") {
			temp := strings.Split(v, ":")

			if len(temp) < 1 {
				continue
			}
			re := regexp.MustCompile(`\b(\d+)`)
			match := re.FindStringSubmatch(v)

			pssValue := 0

			if len(match) > 1 {
				pssValue, _ = strconv.Atoi(match[0])
			}

			if strings.Contains(v, "Java Heap") {
				memInfo.JavaHeap = float64(pssValue) / 1024
			}

			if strings.Contains(v, "Native Heap") {
				memInfo.NativeHeap = float64(pssValue) / 1024
			}

			if strings.Contains(v, "Code") {
				memInfo.Code = float64(pssValue) / 1024
			}

			if strings.Contains(v, "Stack") {
				memInfo.Stack = float64(pssValue) / 1024
			}

			if strings.Contains(v, "Graphics") {
				memInfo.Graphics = float64(pssValue) / 1024
			}

			if strings.Contains(v, "Private Other") {
				memInfo.PrivateOther = float64(pssValue) / 1024
			}

			if strings.Contains(v, "System") {
				memInfo.System = float64(pssValue) / 1024
			}
			if strings.Contains(v, "TOTAL") {
				memInfo.TotalPSS = float64(pssValue) / 1024
			}
		}
	}
	return
}

func getStatusOnPid(device *gadb.Device, pid string) (status *entity.ProcessStatus, err error) {
	data, err1 := device.RunShellCommand(fmt.Sprintf("cat /proc/%s/status", pid))
	if err1 != nil {
		return status, fmt.Errorf("exec command erro : " + fmt.Sprintf("cat /proc/%s/status", pid))
	}
	if strings.Contains(data, "No such file or directory") {
		return nil, errors.New(data)
	}
	scanner := bufio.NewScanner(strings.NewReader(data))
	status = &entity.ProcessStatus{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		var fieldName = strings.TrimRight(fields[0], ":")
		var value = strings.Join(fields[1:], " ")
		switch fieldName {
		case "Name":
			status.Name = value
		case "Umask":
			status.Umask = value
		case "State":
			status.State = value
		case "Tgid":
			status.Tgid = value
		case "Ngid":
			status.Ngid = value
		case "Pid":
			status.Pid = value
		case "PPid":
			status.PPid = value
		case "TracerPid":
			status.TracerPid = value
		case "Uid":
			status.Uid = value
		case "Gid":
			status.Gid = value
		case "FDSize":
			status.FDSize = value
		case "Groups":
			status.Groups = value
		case "VmPeak":
			status.VmPeak = value
		case "VmSize":
			status.VmSize = value
		case "VmLck":
			status.VmLck = value
		case "VmPin":
			status.VmPin = value
		case "VmHWM":
			status.VmHWM = value
		case "VmRSS":
			status.VmRSS = value
		case "RssAnon":
			status.RssAnon = value
		case "RssFile":
			status.RssFile = value
		case "RssShmem":
			status.RssShmem = value
		case "VmData":
			status.VmData = value
		case "VmStk":
			status.VmStk = value
		case "VmExe":
			status.VmExe = value
		case "VmLib":
			status.VmLib = value
		case "VmPTE":
			status.VmPTE = value
		case "VmSwap":
			status.VmSwap = value
		case "Threads":
			status.Threads = value
		case "SigQ":
			status.SigQ = value
		case "SigPnd":
			status.SigPnd = value
		case "ShdPnd":
			status.ShdPnd = value
		case "SigBlk":
			status.SigBlk = value
		case "SigIgn":
			status.SigIgn = value
		case "SigCgt":
			status.SigCgt = value
		case "CapInh":
			status.CapInh = value
		case "CapPrm":
			status.CapPrm = value
		case "CapEff":
			status.CapEff = value
		case "CapBnd":
			status.CapBnd = value
		case "CapAmb":
			status.CapAmb = value
		case "Cpus_allowed":
			status.CpusAllowed = value
		case "Cpus_allowed_list":
			status.CpusAllowedList = value
		case "voluntary_ctxt_switches":
			status.VoluntaryCtxtSwitches = value
		case "nonvoluntary_ctxt_switches":
			status.NonVoluntaryCtxtSwitches = value
		}
	}
	status.TimeStamp = time.Now().UnixMilli()
	return status, err1
}

func newProcessStat(statStr string) (*entity.ProcessStat, error) {
	if strings.Contains(statStr, "No such file or directory") {
		return nil, errors.New(statStr)
	}
	params := strings.Split(statStr, " ")
	var processStat = &entity.ProcessStat{}
	for i, value := range params {
		if i < 24 {
			switch i {
			case 0:
				processStat.Pid = value
			case 1:
				processStat.Comm = value
			case 2:
				processStat.State = value
			case 3:
				processStat.Ppid = value
			case 4:
				processStat.Pgrp = value
			case 5:
				processStat.Session = value
			case 6:
				processStat.Tty_nr = value
			case 7:
				processStat.Tpgid = value
			case 8:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Flags = num
				continue
			case 9:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Minflt = num
				continue
			case 10:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cminflt = num
			case 11:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Majflt = num
			case 12:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Majflt = num
			case 13:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Utime = num
			case 14:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Stime = num
			case 15:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cutime = num
			case 16:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cstime = num
			case 17:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Priority = num
			case 18:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Nice = num
			case 19:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Num_threads = num
			case 20:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Itrealvalue = num
			case 21:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Starttime = num
			case 22:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Vsize = num
			case 23:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Rss = num
			case 24:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Rsslim = num
			}
		}
	}
	processStat.TimeStamp = time.Now().UnixMilli()
	return processStat, nil
}

// https://blog.csdn.net/weixin_39451323/article/details/118083713
func getProcCpuUsage(stat *entity.ProcessStat, pcf *procCpuConfig, nowTotalCPUTime float64) float64 {
	var nowProcCpuTime = float64(stat.Utime) + float64(stat.Stime) + float64(stat.Cutime) + float64(stat.Cstime)
	if pcf.preProcCpuTime == -1.0 {
		pcf.preProcCpuTime = nowProcCpuTime
		pcf.preTotalCpuTime = nowTotalCPUTime
		return 0.0
	}

	cpuUtilization := ((nowProcCpuTime - pcf.preProcCpuTime) / (nowTotalCPUTime - pcf.preTotalCpuTime)) * 100
	pcf.preProcCpuTime = nowProcCpuTime
	pcf.preTotalCpuTime = nowTotalCPUTime

	return cpuUtilization
}

func GetTotalCpuTime(device *gadb.Device) (float64, error) {
	totalCpuInfo, err := device.RunShellCommandWithBytes(fmt.Sprintf("cat /proc/stat"))
	if err != nil {
		return 0, err
	}

	var (
		nowCPU entity.SystemCpuRaw
	)

	//if preCPUMap == nil {
	//	preCPUMap = make(map[string]entity.SystemCpuRaw)
	//}

	scanner := bufio.NewScanner(bytes.NewReader(totalCpuInfo))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") { // changing here if you want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
		}
	}

	totalCpuTime := nowCPU.User + nowCPU.Nice + nowCPU.System + nowCPU.Idle + nowCPU.Iowait + nowCPU.Irq + nowCPU.SoftIrq

	return float64(totalCpuTime), nil
}

var HZ = 100.0 //# ticks/second

func GetProcThreads(device *gadb.Device, perfOption entity.PerfConfig, procThreadCallBackFn func(threadInfo *entity.ProcThreadsInfo, code entity.ServerCode)) {
	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					status, err := getStatusOnPid(device, perfOption.Pid)
					if err != nil {
						log.Error(err)
						procThreadCallBackFn(nil, entity.GetPerfErr)
						return
					}
					var threads int
					if len(status.Threads) > 0 {
						if threads, err = strconv.Atoi(status.Threads); err != nil {
							// todo
							panic(err)
						}
					}
					procThreadCallBackFn(&entity.ProcThreadsInfo{
						Threads:   threads,
						Timestamp: time.Now().UnixMilli(),
					}, entity.RequestSucceed)
				}()
			}
		}
	}()
}

func GetProcCPU(device *gadb.Device, perfOption entity.PerfConfig, procCPUCallBackFn func(cpuInfo *entity.ProcCpuInfo, code entity.ServerCode)) {
	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	pcf := newProcCPUConfig()
	pcf.IntervalTime = perfOption.IntervalTime

	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					stat, err := getStatOnPid(device, perfOption.Pid)
					if err != nil {
						// todo
						log.Error(err)
						procCPUCallBackFn(nil, entity.GetPerfErr)
						return
					}
					nowTotalCPUTime, err := GetTotalCpuTime(device)
					if err != nil {
						log.Error(err)
						procCPUCallBackFn(nil, entity.GetPerfErr)
						return
					}
					procCPUCallBackFn(&entity.ProcCpuInfo{
						CpuUtilization: getProcCpuUsage(stat, pcf, nowTotalCPUTime),
						Timestamp:      time.Now().UnixMilli(),
					}, entity.RequestSucceed)
				}()
			}
		}
	}()
}

func GetProcMem(device *gadb.Device, perfOption entity.PerfConfig, procMemCallBackFn func(memInfo *entity.ProcMemInfo, code entity.ServerCode)) {
	timer := time.Tick(time.Duration(perfOption.IntervalTime * int(time.Second)))
	pcf := newProcCPUConfig()
	pcf.IntervalTime = perfOption.IntervalTime
	go func() {
		for {
			select {
			case <-perfOption.Ctx.Done():
				return
			case <-timer:
				go func() {
					mem, err := getMemTotalPSS(device, perfOption.Pid)
					if err != nil {
						log.Error(err)
						procMemCallBackFn(nil, entity.GetPerfErr)
						return
					}
					procMemCallBackFn(mem, entity.RequestSucceed)
				}()
			}
		}
	}()
}
