package entity

type ProcessIO struct {
	Rchar               int `json:"rchar"`
	Wchar               int `json:"wchar"`
	Syscr               int `json:"syscr"`
	Syscw               int `json:"syscw"`
	ReadBytes           int `json:"readBytes"`
	WriteBytes          int `json:"writeBytes"`
	CancelledWriteBytes int `json:"cancelledWriteBytes"`
}

type ProcessStat struct {
	Pid         string
	Comm        string
	State       string
	Ppid        string
	Pgrp        string
	Session     string
	Tty_nr      string
	Tpgid       string
	Flags       int
	Minflt      int
	Cminflt     int
	Majflt      int
	Cmajflt     int
	Utime       int
	Stime       int
	Cutime      int
	Cstime      int
	Priority    int
	Nice        int
	Num_threads int
	Itrealvalue int
	Starttime   int
	Vsize       int
	Rss         int
	Rsslim      int
	TimeStamp   int64
}

type ProcessStatus struct {
	Name                     string `json:"name"`
	Umask                    string `json:"umask"`
	State                    string `json:"state"`
	Tgid                     string `json:"tgid"`
	Ngid                     string `json:"ngid"`
	Pid                      string `json:"pid"`
	PPid                     string `json:"pPid"`
	TracerPid                string `json:"tracerPid"`
	Uid                      string `json:"uid"`
	Gid                      string `json:"gid"`
	FDSize                   string `json:"fdSize"`
	Groups                   string `json:"groups"`
	VmPeak                   string `json:"vmPeak"`
	VmSize                   string `json:"vmSize"`
	VmLck                    string `json:"vmLck"`
	VmPin                    string `json:"vmPin"`
	VmHWM                    string `json:"vmHWM"`
	VmRSS                    string `json:"vmRSS"`
	RssAnon                  string `json:"rssAnon"`
	RssFile                  string `json:"rssFile"`
	RssShmem                 string `json:"rssShmem"`
	VmData                   string `json:"vmData"`
	VmStk                    string `json:"vmStk"`
	VmExe                    string `json:"vmExe"`
	VmLib                    string `json:"vmLib"`
	VmPTE                    string `json:"vmPTE"`
	VmSwap                   string `json:"vmSwap"`
	Threads                  string `json:"threads"`
	SigQ                     string `json:"sigQ"`
	SigPnd                   string `json:"sigPnd"`
	ShdPnd                   string `json:"shdPnd"`
	SigBlk                   string `json:"sigBlk"`
	SigIgn                   string `json:"sigIgn"`
	SigCgt                   string `json:"sigCgt"`
	CapInh                   string `json:"capInh"`
	CapPrm                   string `json:"capPrm"`
	CapEff                   string `json:"capEff"`
	CapBnd                   string `json:"capBnd"`
	CapAmb                   string `json:"capAmb"`
	CpusAllowed              string `json:"cpusAllowed"`
	CpusAllowedList          string `json:"cpusAllowedList"`
	VoluntaryCtxtSwitches    string `json:"voluntaryCtxtSwitches"`
	NonVoluntaryCtxtSwitches string `json:"nonVoluntaryCtxtSwitches"`
	TimeStamp                int64
}

type ProcessInfo struct {
	CPUInfo    *ProcCpuInfo     `json:"cpuInfo,omitempty"`
	MemInfo    *ProcMemInfo     `json:"memInfo,omitempty"`
	ThreadInfo *ProcThreadsInfo `json:"threadInfo,omitempty"`
	Error      []string         `json:"error,omitempty"`
}

// 都是MB为单位
type ProcMemInfo struct {
	UUID         string  `json:"-" gorm:"primaryKey"`
	TotalPSS     float64 `json:"totalPSS" xlsx:"A-totalPSS"`
	JavaHeap     float64 `json:"javaHeap" xlsx:"B-javaHeap"`
	NativeHeap   float64 `json:"nativeHeap" xlsx:"C-nativeHeap"`
	Code         float64 `json:"code" xlsx:"D-code"`
	Stack        float64 `json:"stack" xlsx:"E-stack"`
	Graphics     float64 `json:"graphics" xlsx:"F-graphics"`
	PrivateOther float64 `json:"privateOther" xlsx:"G-privateOther"`
	System       float64 `json:"system" xlsx:"H-system"`
	//PhyRSS    int   `json:"phyRSS"`
	//VmSize    int   `json:"vmRSS"`
	Timestamp int64 `json:"timestamp" gorm:"primaryKey" xlsx:"I-timestamp"`
}

type ProcCpuInfo struct {
	UUID           string  `json:"-" gorm:"primaryKey"`
	CpuUtilization float64 `json:"cpuUtilization"  xlsx:"A-cpuUtilization"` // 百分比
	Timestamp      int64   `json:"timestamp"  xlsx:"B-timestamp" gorm:"primaryKey"`
}

type ProcThreadsInfo struct {
	UUID      string `json:"-" gorm:"primaryKey" `
	Threads   int    `json:"threadCount"   xlsx:"A-threadCount"`
	Timestamp int64  `json:"timestamp"   xlsx:"B-timestamp" gorm:"primaryKey"`
}
