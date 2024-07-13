package perf_test

import (
	"context"
	"fionna/android/android_util"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"testing"
	"time"
)

func TestPerfTool(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "emulator-5554")
	if err != nil {
		panic(err)
	}

	ctx, exitFn := context.WithCancel(context.Background())
	framePerf := perf.NewPerfTool(device, ctx)

	framePerf.Init()
	framePerf.GetFrame(func(frame *entity.SysFrameInfo, code entity.ServerCode) {
		fmt.Println(frame)
	})
	time.Sleep(100 * time.Second)
	exitFn()

}
