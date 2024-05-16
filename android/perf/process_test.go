package perf_test

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"testing"
	"time"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestGetProcMem(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	pid, err := android_util.GetPidOnPackageName(device, "com.taou.maimai")
	if err != nil {
		panic(err)
	}
	perfOption := entity.PerfConfig{Pid: pid, IntervalTime: 1}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption.Ctx = ctx

	perf.GetProcMem(device, perfOption, func(memInfo *entity.ProcMemInfo, code entity.ServerCode) {
		data, _ := json.Marshal(memInfo)
		fmt.Println(string(data))
	})

	time.Sleep(10 * time.Second)
	exitFn()
}

func TestGetProcThreads(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	pid, err := android_util.GetPidOnPackageName(device, "com.taou.maimai")
	if err != nil {
		panic(err)
	}
	perfOption := entity.PerfConfig{Pid: pid, IntervalTime: 1}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption.Ctx = ctx

	perf.GetProcThreads(device, perfOption, func(threadInfo *entity.ProcThreadsInfo, code entity.ServerCode) {
		data, _ := json.Marshal(threadInfo)
		fmt.Println(string(data))
	})

	time.Sleep(20 * time.Second)
	exitFn()
}

func TestGetProcCpu(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	pid, err := android_util.GetPidOnPackageName(device, "com.baidu.tieba_mini")
	if err != nil {
		panic(err)
	}
	perfOption := entity.PerfConfig{Pid: pid, IntervalTime: 1}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption.Ctx = ctx

	perf.GetProcCPU(device, perfOption, func(cpuInfo *entity.ProcCpuInfo, code entity.ServerCode) {
		data, _ := json.Marshal(cpuInfo)
		fmt.Println(string(data))
	})

	time.Sleep(20 * time.Second)
	exitFn()
}
