package entity

type SystemNetworkData struct {
	UUID      string `json:"-" gorm:"primaryKey"`
	Data      string `json:"data" gorm:"primaryKey"    xlsx:"A-data"`
	Timestamp int64  `json:"timestamp,omitempty" gorm:"primaryKey"   xlsx:"B-timestamp"`
}

type SystemNetworkInfo struct {
	InterfaceName string  `json:"interfaceName"  gorm:"primaryKey"    xlsx:"A-interfaceName"`
	UUID          string  `json:"-" gorm:"primaryKey"`
	Timestamp     int64   `json:"timestamp,omitempty" gorm:"primaryKey"   xlsx:"B-timestamp"`
	IPv4          string  `json:"ipv4"   xlsx:"C-ipv4"`
	IPv6          string  `json:"ipv6" xlsx:"D-ipv6"`
	Rx            float64 `json:"rx"  xlsx:"E-rx"` // MB
	Tx            float64 `json:"tx"  xlsx:"F-tx"` // MB
}

type SystemCpuRaw struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

type SystemCPUData struct {
	UUID      string `json:"-" gorm:"primaryKey"`
	Data      string `json:"data" gorm:"primaryKey"    xlsx:"A-data"`
	Timestamp int64  `json:"timestamp,omitempty" gorm:"primaryKey"   xlsx:"B-timestamp"`
}

type SystemCPUInfo struct {
	CPUName   string  `json:"cpuName" gorm:"primaryKey"   xlsx:"A-cpuName"`
	UUID      string  `json:"-" gorm:"primaryKey"`
	User      float32 `json:"user" gorm:"-"`
	Nice      float32 `json:"nice" gorm:"-"`
	System    float32 `json:"system" gorm:"-"`
	Idle      float32 `json:"idle" gorm:"-"`
	Iowait    float32 `json:"iowait" gorm:"-"`
	Irq       float32 `json:"irq" gorm:"-"`
	SoftIrq   float32 `json:"softIrq" gorm:"-"`
	Steal     float32 `json:"steal" gorm:"-"`
	Guest     float32 `json:"guest" gorm:"-"`
	Usage     float32 `json:"cpuUsage"   xlsx:"B-cpuUsage"` // 百分比
	Timestamp int64   `json:"timestamp" gorm:"primaryKey"   xlsx:"C-timestamp"`
}

type SystemMemInfo struct {
	UUID       string `json:"-" gorm:"primaryKey"`
	MemTotal   uint64 `json:"memTotal"   xlsx:"A-memTotal"`
	MemFree    uint64 `json:"memFree" gorm:"-"`
	MemBuffers uint64 `json:"memBuffers" gorm:"-"`
	MemCached  uint64 `json:"memCached" gorm:"-"`
	MemUsage   uint64 `json:"memUsage"   xlsx:"B-memUsage"`
	SwapTotal  uint64 `json:"swapTotal" gorm:"-"`
	SwapFree   uint64 `json:"swapFree" gorm:"-"`
	Timestamp  int64  `json:"timestamp" gorm:"primaryKey"   xlsx:"C-timestamp"`
}

// MB
type SystemInfo struct {
	MemInfo     *SystemMemInfo                `json:"memInfo,omitempty"`
	NetworkInfo map[string]*SystemNetworkInfo `json:"networkInfo,omitempty"`
	CPU         map[string]*SystemCPUInfo     `json:"cpuInfo,omitempty"`
	Frame       *SysFrameInfo                 `json:"frame,omitempty"`
	Temperature *SysTemperature               `json:"temperature,omitempty"`
	Error       []string                      `json:"error,omitempty"`
	//Timestamp   int64                         `json:"timeStamp"`
}

type SysFrameInfo struct {
	UUID         string `json:"-" gorm:"primaryKey"`
	Timestamp    int64  `json:"timestamp,omitempty" gorm:"primaryKey"  xlsx:"A-timestamp"`
	FPS          int    `json:"FPS"  xlsx:"B-FPS"`
	JankCount    int    `json:"jankCount"  xlsx:"C-jankCount"`
	BigJankCount int    `json:"bigJankCount"  xlsx:"D-bigJankCount"`
}

type SysTemperature struct {
	Temperature float64 `json:"temperature,omitempty"  xlsx:"A-temperature"`
	UUID        string  `json:"-" gorm:"primaryKey"`
	Timestamp   int64   `json:"timestamp,omitempty" gorm:"primaryKey"  xlsx:"B-timestamp"`
}
