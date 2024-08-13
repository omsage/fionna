package cmd

import (
	"context"
	"fionna/android/android_util"
	"fionna/android/gadb"
	"fionna/entity"
	"fionna/server/android"
	"fionna/server/db"
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
		db.InitDB(dbName)
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

		serialInfo := android_util.GetSerialInfo(device)

		android.InitPerfAndStart(&serialInfo, perfConfig, device, nil)

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
	cliPerfCmd.Flags().BoolVar(&perfConfig.SysTemperature, "sys-temperature", false, "get system temperature data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.FPS, "fps", false, "get fps data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.Jank, "jank", false, "get jank data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcThread, "proc-threads", false, "get process threads")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcCpu, "proc-cpu", false, "get process cpu data")
	cliPerfCmd.Flags().BoolVar(&perfConfig.ProcMem, "proc-mem", false, "get process mem data")
	cliPerfCmd.Flags().StringVar(&dbName, "db-path", "test.db", "specify the SQLite path to use")
}
