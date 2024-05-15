package entity

type OverallSummary struct {
	NetworkSummary        map[string]SystemNetworkSummary `json:"networkSummary,omitempty"  xlsx:"A-networkSummary	"`
	SysCpuSummary         map[string]SystemCPUSummary     `json:"sysCpuSummary,omitempty"  xlsx:"B-sysCpuSummary"`
	SysMemSummary         *SystemMemSummary               `json:"sysMemSummary,omitempty"  xlsx:"C-sysMemSummary"`
	SysFrameSummary       *FrameSummary                   `json:"sysFrameSummary,omitempty"  xlsx:"D-sysFrameSummary"`
	SysTemperatureSummary *SystemTemperatureSummary       `json:"sysTemperatureSummary,omitempty"  xlsx:"E-sysTemperatureSummary"`
	ProcCpu               *ProcCpuSummary                 `json:"procCpuSummary,omitempty"  xlsx:"F-procCpuSummary"`
	ProcMem               *ProcMemSummary                 `json:"procMemSummary,omitempty"  xlsx:"G-procMemSummary"`
}

func NewSystemNetworkSummary(uuid string, name string) *SystemNetworkSummary {
	return &SystemNetworkSummary{
		UUID: uuid,
		Name: name,
	}
}

type SystemNetworkSummary struct {
	UUID         string  `json:"-" gorm:"primaryKey" `
	Name         string  `json:"name" gorm:"primaryKey"  xlsx:"A-name"`
	AllSysTxData float64 `json:"allSysTxData"  xlsx:"C-allSysTxData"`
	AllSysRxData float64 `json:"allSysRxData"  xlsx:"B-allSysRxData"`
}

func NewSystemTemperature(uuid string) *SystemTemperatureSummary {
	return &SystemTemperatureSummary{
		UUID:           uuid,
		MaxTemperature: -1,
	}
}

type SystemTemperatureSummary struct {
	UUID            string  `json:"-" gorm:"primaryKey"  `
	DiffTemperature float64 `json:"diffTemperature"  xlsx:"B-diffTemperature"`
	MaxTemperature  float64 `json:"mxTemperature"  xlsx:"A-mxTemperature"`
}

func NewSystemCPUSummary(uuid string) *SystemCPUSummary {
	return &SystemCPUSummary{
		UUID:      uuid,
		MaxSysCPU: -1,
	}
}

type SystemCPUSummary struct {
	UUID      string  `json:"-" gorm:"primaryKey"  `
	CpuName   string  `json:"cpuName" gorm:"primaryKey"  xlsx:"A-cpuName"`
	AvgSysCPU float64 `json:"avgSysCPU"  xlsx:"B-avgSysCPU"`
	MaxSysCPU float64 `json:"maxSysCPU"  xlsx:"C-maxSysCPU"`
}

func NewSystemMemSummary(uuid string) *SystemMemSummary {
	return &SystemMemSummary{
		UUID:        uuid,
		MaxMemTotal: -1,
	}
}

type SystemMemSummary struct {
	UUID        string  `json:"-" gorm:"primaryKey"  `
	MaxMemTotal float64 `json:"maxMemTotal"  xlsx:"A-maxMemTotal"`
	AvgMemTotal float64 `json:"avgMemTotal"  xlsx:"B-avgMemTotal"`
}

func NewFrameSummary(uuid string) *FrameSummary {
	return &FrameSummary{
		UUID:            uuid,
		MaxJankCount:    -1,
		MaxBigJankCount: -1,
	}
}

type FrameSummary struct {
	UUID             string  `json:"-" gorm:"primaryKey"`
	AvgFPS           float64 `json:"avgFPS"  xlsx:"A-avgFPS"`
	AllJankCount     int     `json:"allJankCount"  xlsx:"B-allJankCount"`
	AllBigJankCount  int     `json:"allBigJankCount"  xlsx:"C-allBigJankCount"`
	MaxJankCount     int     `json:"maxJankCount"  xlsx:"D-maxJankCount"`
	MaxBigJankCount  int     `json:"maxBigJankCount"  xlsx:"E-maxBigJankCount"`
	JankCountRate    float64 `json:"jankCountRate"  xlsx:"F-jankCountRate"`       // 百分比
	BigJankCountRate float64 `json:"bigJankCountRate"  xlsx:"G-bigJankCountRate"` // 百分比
}

func NewProcCpuSummary(uuid string) *ProcCpuSummary {
	return &ProcCpuSummary{
		UUID:       uuid,
		MaxProcCPU: -1,
	}
}

type ProcCpuSummary struct {
	UUID       string  `json:"-" gorm:"primaryKey"`
	AvgProcCPU float64 `json:"avgProcCPU" gorm:"avgProcCPU"  xlsx:"A-avgProcCPU"`
	MaxProcCPU float64 `json:"maxProcCPU" gorm:"maxProcCPU"  xlsx:"B-maxProcCPU"`
}

func NewProcMemSummary(uuid string) *ProcMemSummary {
	return &ProcMemSummary{
		UUID:            uuid,
		MaxTotalPSS:     -1,
		MaxJavaHeap:     -1,
		MaxNativeHeap:   -1,
		MaxCode:         -1,
		MaxStack:        -1,
		MaxGraphics:     -1,
		MaxPrivateOther: -1,
		MaxSystem:       -1,
	}
}

type ProcMemSummary struct {
	UUID            string  `json:"-" gorm:"primaryKey" `
	AvgTotalPSS     float64 `json:"avgTotalPSS"  xlsx:"A-avgTotalPSS"`
	AvgJavaHeap     float64 `json:"avgJavaHeap"  xlsx:"B-avgJavaHeap"`
	AvgNativeHeap   float64 `json:"avgNativeHeap"  xlsx:"C-avgNativeHeap"`
	AvgCode         float64 `json:"avgCode"  xlsx:"D-avgCode"`
	AvgStack        float64 `json:"avgStack"  xlsx:"E-avgStack"`
	AvgGraphics     float64 `json:"avgGraphics"  xlsx:"F-avgGraphics"`
	AvgPrivateOther float64 `json:"avgPrivateOther" xlsx:"G-avgPrivateOther"`
	AvgSystem       float64 `json:"avgSystem"  xlsx:"H-avgSystem"`
	MaxTotalPSS     float64 `json:"maxTotalPSS"  xlsx:"I-maxTotalPSS"`
	MaxJavaHeap     float64 `json:"maxJavaHeap"  xlsx:"J-maxJavaHeap"`
	MaxNativeHeap   float64 `json:"maxNativeHeap"  xlsx:"K-maxNativeHeap"`
	MaxCode         float64 `json:"maxCode"  xlsx:"L-maxCode"`
	MaxStack        float64 `json:"maxStack"  xlsx:"M-maxStack"`
	MaxGraphics     float64 `json:"maxGraphics"  xlsx:"N-maxGraphics"`
	MaxPrivateOther float64 `json:"maxPrivateOther"  xlsx:"O-maxPrivateOther"`
	MaxSystem       float64 `json:"maxSystem"  xlsx:"P-maxSystem"`
}
