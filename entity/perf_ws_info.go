package entity

import "context"

type PerfRecvMessageType string
type PerfSendMessageType string

const (
	StartPerfType PerfRecvMessageType = "startPerfmon"
	ClosePerfType PerfRecvMessageType = "closePerfmon"
	PongPerfType  PerfRecvMessageType = "pongPerfmon"
	PerfDataType  PerfSendMessageType = "perfdata"
	PerfErrorType PerfSendMessageType = "error"
)

type PerfRecvMessage struct {
	MessageType PerfRecvMessageType `json:"messageType"`
	Data        interface{}         `json:"data"`
}

type PerfData struct {
	SystemPerfData SystemInfo  `json:"system"`
	ProcPerfData   ProcessInfo `json:"process"`
}

func NewPerfDataMessage(PerfData *PerfData) *PerfDataMessage {
	return &PerfDataMessage{
		MessageType: PerfDataType,
		Data:        PerfData,
	}
}

type PerfDataMessage struct {
	MessageType PerfSendMessageType `json:"messageType"`
	Data        interface{}         `json:"perfData"`
}

func NewPerfDataError(message string) *PerfErrorMessage {
	return &PerfErrorMessage{
		MessageType: PerfErrorType,
		ErrorInfo:   message,
	}
}

type PerfErrorMessage struct {
	MessageType PerfSendMessageType `json:"messageType"`
	ErrorInfo   string              `json:"errorInfo"`
}

type PerfConfig struct {
	SysCpu         bool `json:"sysCpu" gorm:"sysCpu"`
	SysMem         bool `json:"sysMem"  gorm:"sysMem"`
	SysNetwork     bool `json:"sysNetwork"  gorm:"sysNetwork"`
	SysTemperature bool `json:"sysTemperature"`
	FPS            bool `json:"FPS"  gorm:"FPS"`
	Jank           bool `json:"jank"  gorm:"jank"`
	ProcCpu        bool `json:"procCpu"  gorm:"procCpu"`
	ProcMem        bool `json:"procMem"  gorm:"procMem"`
	ProcThread     bool `json:"procThread"  gorm:"procThread"`

	UUID         string             `json:"uuid,omitempty" gorm:"primaryKey"`
	Serial       string             `json:"serial,omitempty" gorm:"serial"`
	Pid          string             `json:"pid,omitempty"  gorm:"-"`
	PackageName  string             `json:"packageName,omitempty"  gorm:"packageName"`
	IntervalTime int                `json:"intervalTime,omitempty"   gorm:"intervalTime"`
	Ctx          context.Context    `json:"-"   gorm:"-"`
	CancelFn     context.CancelFunc `json:"-"   gorm:"-"`
}

func NewPerfOption(ctx context.Context, IntervalTime int, opts ...PerfOption) *PerfConfig {
	pcf := &PerfConfig{
		Ctx:          ctx,
		IntervalTime: IntervalTime,
	}
	for _, opt := range opts {
		opt(pcf)
	}
	return pcf
}

type PerfOption func(*PerfConfig)

func WithPid(pid string) PerfOption {
	return func(option *PerfConfig) {
		option.Pid = pid
	}
}
