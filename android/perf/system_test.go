package perf_test

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"testing"
	"time"
)

func TestGetSysMem(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())

	perfOption := entity.NewPerfOption(ctx, 1)

	perf.GetSysMem(device, *perfOption, func(sysMem *entity.SystemMemInfo, code entity.ServerCode) {
		data, _ := json.Marshal(sysMem)
		fmt.Println(string(data))
	})

	time.Sleep(10 * time.Second)
	exitFn()
}

func TestGetSysFrame(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption := entity.NewPerfOption(ctx, 1)

	perf.GetSysFrame(device, *perfOption, func(frame *entity.SysFrameInfo, code entity.ServerCode) {
		data, _ := json.Marshal(frame)
		fmt.Println(string(data))
	})

	time.Sleep(20 * time.Second)
	exitFn()
}

func TestGetSysNetwork(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption := entity.NewPerfOption(ctx, 1)

	perf.GetSysNetwork(device, *perfOption, func(sysNet map[string]*entity.SystemNetworkInfo, code entity.ServerCode) {
		data, _ := json.Marshal(sysNet)
		fmt.Println(string(data))
	})

	time.Sleep(20 * time.Second)
	exitFn()
}

func TestGetSysCpu(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption := entity.NewPerfOption(ctx, 1)

	perf.GetSysCPU(device, *perfOption, func(CPU map[string]*entity.SystemCPUInfo, code entity.ServerCode) {
		data, _ := json.Marshal(CPU)
		fmt.Println(string(data))
	})

	//perf.GetSysCPU()

	time.Sleep(20 * time.Second)
	exitFn()
}

func TestGetSysTemperature(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())
	perfOption := entity.NewPerfOption(ctx, 1)

	perf.GetSysTemperature(device, *perfOption, func(info *entity.SysTemperature, code entity.ServerCode) {
		data, _ := json.Marshal(info)
		fmt.Println(string(data))
	})

	//perf.GetSysCPU()

	time.Sleep(60 * time.Second)
	exitFn()
}
