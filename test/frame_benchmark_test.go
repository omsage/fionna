package test

import (
	"bufio"
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

type RenderTime struct {
	Draw    float64
	Prepare float64
	Process float64
	Execute float64
}

func sum(arr []float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	return sum(arr, n-1) + arr[n-1]
}

func NewFrameBenchmark(device *gadb.Device, pkg string) *FrameBenchmark {
	ctx, exitFn := context.WithCancel(context.Background())
	pid, err := android_util.GetPidOnPackageName(device, pkg)
	if err != nil {
		panic(err)
	}
	return &FrameBenchmark{
		ctx:    ctx,
		exitFn: exitFn,
		device: device,
		pkg:    pkg,
		pid:    pid,
	}
}

type FrameBenchmark struct {
	ctx    context.Context
	exitFn context.CancelFunc
	device *gadb.Device
	pkg    string
	pid    string
}

func (f *FrameBenchmark) getProcessFPSBySurfaceFlinger() *entity.SysFrameInfo {
	_, err := f.device.RunShellCommand("dumpsys SurfaceFlinger --latency-clear")
	lines, err := f.device.RunShellLoopCommand(
		fmt.Sprintf("dumpsys SurfaceFlinger | grep %s", f.pkg))
	if err != nil {
		panic(err)
	}

	activity := ""

	scanner := bufio.NewScanner(lines)
	reg := regexp.MustCompile("\\[.*#0|\\(.*\\)")
	for scanner.Scan() {
		line := scanner.Text()

		activity = reg.FindString(line)

		if activity == "" {
			continue
		}
		break
	}
	if activity == "" {
		panic(fmt.Sprintf("could not find app %s activity", f.pkg))
	}
	r := strings.NewReplacer("[", "", "(", "", ")", "")
	activity = r.Replace(activity)

	lines, err = f.device.RunShellLoopCommand(
		fmt.Sprintf("dumpsys SurfaceFlinger --latency '%s'", activity))
	if err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(lines)
	var preFrame float64
	var t []float64
	for scanner.Scan() {
		line := scanner.Text()
		l := strings.Split(line, "\t")
		if len(l) < 3 {
			continue
		}
		if l[0][0] == '0' {
			continue
		}
		frame, _ := strconv.ParseFloat(l[1], 64)
		if frame == math.MaxInt64 {
			continue
		}
		frame /= 1e6
		if frame <= preFrame {
			continue
		}
		if preFrame == 0 {
			preFrame = frame
			continue
		}
		t = append(t, frame-preFrame)
		preFrame = frame
	}

	le := len(t)
	if le == 0 {
		return &entity.SysFrameInfo{Timestamp: time.Now().UnixMilli()}
	}

	frame := &entity.SysFrameInfo{
		FPS:       (int)(float64(le) * 1000 / (sum(t, le))),
		Timestamp: time.Now().UnixMilli(),
	}
	return frame
}

func (f *FrameBenchmark) getProcessFPSByGFXInfo() *entity.SysFrameInfo {
	lines, err := f.device.RunShellLoopCommand(
		fmt.Sprintf("dumpsys gfxinfo %s | grep '.*visibility=0' -A129 | grep Draw -A128 | grep 'View hierarchy:' -B129", f.pid))
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(lines)
	frameCount := 0
	vsyncCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Draw") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		frameCount++
		s := strings.Split(line, "\t")
		if len(s) == 5 {
			render := RenderTime{}
			render.Draw, _ = strconv.ParseFloat(s[1], 64)
			render.Prepare, _ = strconv.ParseFloat(s[2], 64)
			render.Process, _ = strconv.ParseFloat(s[3], 64)
			render.Execute, _ = strconv.ParseFloat(s[4], 64)
			total := render.Draw + render.Prepare + render.Process + render.Execute

			if total > 16.67 {
				vsyncCount += (int)(math.Ceil(total/16.67) - 1)
			}
		}
	}
	if (frameCount + vsyncCount) == 0 {
		return &entity.SysFrameInfo{
			Timestamp: time.Now().UnixMilli(),
		}
	}
	frame := &entity.SysFrameInfo{
		FPS:       frameCount * 60 / (frameCount + vsyncCount),
		Timestamp: time.Now().UnixMilli(),
	}
	return frame
}

func (f *FrameBenchmark) GetFPSBySurfaceFlinger(frameCallback func(frame *entity.SysFrameInfo, code entity.ServerCode)) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			frame := f.getProcessFPSBySurfaceFlinger()
			frameCallback(frame, entity.RequestSucceed)
		case <-f.ctx.Done():
			break
		}
	}
}

func (f *FrameBenchmark) GetFPSByGFXInfo(frameCallback func(frame *entity.SysFrameInfo, code entity.ServerCode)) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			frame := f.getProcessFPSByGFXInfo()
			frameCallback(frame, entity.RequestSucceed)
		case <-f.ctx.Done():
			break
		}
	}
}

func (f *FrameBenchmark) GetFPSByFrameTool(frameCallback func(frame *entity.SysFrameInfo, code entity.ServerCode)) {
	framePerf := perf.NewPerfTool(f.device, f.ctx)

	framePerf.Init()

	framePerf.GetFrame(frameCallback)
}

func (f *FrameBenchmark) Exit() {
	f.exitFn()
}

func TestFPSBySurfaceFlinger(t *testing.T) {
	SetClient()
	testPkg := "com.dragonli.projectsnow.lhm"
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	frameBenchmark := NewFrameBenchmark(device, testPkg)
	go func() {
		frameBenchmark.GetFPSBySurfaceFlinger(func(frame *entity.SysFrameInfo, code entity.ServerCode) {
			data, _ := json.Marshal(frame)
			fmt.Println(string(data))
		})
	}()

	time.Sleep(10 * time.Second)
	frameBenchmark.Exit()
}

func TestFPSByGFXInfo(t *testing.T) {
	SetClient()
	testPkg := "com.dragonli.projectsnow.lhm"
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	frameBenchmark := NewFrameBenchmark(device, testPkg)
	go func() {
		frameBenchmark.GetFPSByGFXInfo(func(frame *entity.SysFrameInfo, code entity.ServerCode) {
			data, _ := json.Marshal(frame)
			fmt.Println(string(data))
		})
	}()
	time.Sleep(10 * time.Second)
	frameBenchmark.Exit()
}

func TestFPSByOMSageFrameTool(t *testing.T) {
	SetClient()
	testPkg := "com.dragonli.projectsnow.lhm"
	device, err := android_util.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	frameBenchmark := NewFrameBenchmark(device, testPkg)
	go func() {
		frameBenchmark.GetFPSByFrameTool(func(frame *entity.SysFrameInfo, code entity.ServerCode) {
			data, _ := json.Marshal(frame)
			fmt.Println(string(data))
		})
	}()
	time.Sleep(10 * time.Second)
	frameBenchmark.Exit()
}
