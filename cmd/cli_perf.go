package cmd

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/android/perf"
	"fionna/entity"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var cliPerfCmd = &cobra.Command{
	Use:   "cli-perf",
	Short: "Fionna perf cli mode",
	Long:  "Fionna perf cli mode",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		client, err := gadb.NewClient()
		if err != nil {
			panic(err)
		}
		device, err := android_util.GetDevice(client, serial)

		if err != nil {
			panic(err)
		}

		if (perfConfig.Pid == "" && perfConfig.PackageName == "") &&
			!perfConfig.SysCpu &&
			!perfConfig.SysMem &&
			!perfConfig.SysNetwork &&
			!perfConfig.SysTemperature &&
			!perfConfig.FPS &&
			!perfConfig.Jank {
			sysAllParamsSet()
		}
		if (perfConfig.Pid != "" || perfConfig.PackageName != "") &&
			!perfConfig.ProcMem &&
			!perfConfig.ProcCpu &&
			!perfConfig.ProcThread {
			perfConfig.ProcMem = true
			perfConfig.ProcCpu = true
			perfConfig.ProcThread = true
		}

		if perfConfig.PackageName != "" {
			perfConfig.Pid, err = android_util.GetPidOnPackageName(device, perfConfig.PackageName)
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)

		exitCtx, exitChancel := context.WithCancel(context.Background())

		perfConfig.Ctx = exitCtx
		perfConfig.CancelFn = exitChancel

		startGetPerf(device, *perfConfig)
		<-sig
		exitChancel()
		os.Exit(0)
		return nil
	},
}

func sysAllParamsSet() {
	perfConfig.SysCpu = true
	perfConfig.SysTemperature = true
	perfConfig.SysMem = true
	perfConfig.SysNetwork = true
	perfConfig.FPS = true
	perfConfig.Jank = true
}

func startGetPerf(device *gadb.Device, config entity.PerfConfig) {
	if config.FPS || config.Jank {

		go func() {
			perf.GetSysFrame(device, config, func(frame *entity.SysFrameInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{SystemPerfData: &entity.SystemInfo{Frame: frame}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}
	if config.SysCpu {

		go func() {
			perf.GetSysCPU(device, config, func(CPU map[string]*entity.SystemCPUInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{SystemPerfData: &entity.SystemInfo{CPU: CPU}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}
	if config.SysMem {

		go func() {
			perf.GetSysMem(device, config, func(sysMem *entity.SystemMemInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{SystemPerfData: &entity.SystemInfo{MemInfo: sysMem}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}
	if config.SysNetwork {
		go func() {
			perf.GetSysNetwork(device, config, func(sysNet map[string]*entity.SystemNetworkInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{SystemPerfData: &entity.SystemInfo{NetworkInfo: sysNet}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}

	if config.ProcCpu {

		go func() {
			perf.GetProcCPU(device, config, func(cpuInfo *entity.ProcCpuInfo, code entity.ServerCode) {

				data, err := json.Marshal(&entity.PerfData{ProcPerfData: &entity.ProcessInfo{CPUInfo: cpuInfo}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}

	if config.ProcMem {

		go func() {
			perf.GetProcMem(device, config, func(memInfo *entity.ProcMemInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{ProcPerfData: &entity.ProcessInfo{MemInfo: memInfo}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()

	}

	if config.ProcThread {
		go func() {
			perf.GetProcThreads(device, config, func(threadInfo *entity.ProcThreadsInfo, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{ProcPerfData: &entity.ProcessInfo{ThreadInfo: threadInfo}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()
	}

	if config.SysTemperature {
		go func() {
			perf.GetSysTemperature(device, config, func(temperatureInfo *entity.SysTemperature, code entity.ServerCode) {
				data, err := json.Marshal(&entity.PerfData{SystemPerfData: &entity.SystemInfo{Temperature: temperatureInfo}})
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			})
		}()
	}
}

var (
	serial     string
	perfConfig = &entity.PerfConfig{
		IntervalTime: 1,
	}
)

func init() {
	rootCmd.AddCommand(cliPerfCmd)
	cliPerfCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial (default first device)")
	cliPerfCmd.Flags().StringVarP(&perfConfig.Pid, "pid", "d", "", "get PID data")
	cliPerfCmd.Flags().StringVarP(&perfConfig.PackageName, "package", "p", "", "app package name")
	cliPerfCmd.Flags().BoolVar(&perfConfig.SysCpu, "sys-cpu", false, "get system cpu data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.SysMem, "sys-mem", false, "get system memory data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.SysNetwork, "sys-network", false, "get system networking data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.FPS, "fps", false, "get fps data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.Jank, "jank", false, "get jank data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcThread, "proc-threads", false, "get process threads")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcCpu, "proc-cpu", false, "get process cpu data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcMem, "proc-mem", false, "get process mem data")
}
