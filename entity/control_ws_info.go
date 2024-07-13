package entity

import "context"

type ControlRecvMessageType string
type ControlSendMessageType string

const (
	ControlTouchType ControlRecvMessageType = "touch"
	StartPerfType    ControlRecvMessageType = "startPerfmon"
	ClosePerfType    ControlRecvMessageType = "closePerfmon"
	PongPerfType     ControlRecvMessageType = "pongPerfmon"
	PerfDataType     ControlSendMessageType = "perfdata"
	RotationDataType ControlSendMessageType = "rotation"
	PerfErrorType    ControlSendMessageType = "perfError"
)

type PerfRecvMessage struct {
	MessageType ControlRecvMessageType `json:"messageType"`
	Data        interface{}            `json:"data"`
}

type PerfData struct {
	SystemPerfData *SystemInfo  `json:"system,omitempty"`
	ProcPerfData   *ProcessInfo `json:"process,omitempty"`
}

func NewPerfDataMessage(PerfData *PerfData) *PerfDataMessage {
	return &PerfDataMessage{
		MessageType: PerfDataType,
		Data:        PerfData,
	}
}

type PerfDataMessage struct {
	MessageType ControlSendMessageType `json:"messageType"`
	Data        interface{}            `json:"perfData"`
}

func NewPerfDataError(message string) *PerfErrorMessage {
	return &PerfErrorMessage{
		MessageType: PerfErrorType,
		ErrorInfo:   message,
	}
}

type PerfErrorMessage struct {
	MessageType ControlSendMessageType `json:"messageType"`
	ErrorInfo   string                 `json:"errorInfo"`
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

func NewPerfOption(ctx context.Context, opts ...PerfOption) *PerfConfig {
	pcf := &PerfConfig{
		Ctx:          ctx,
		IntervalTime: 1,
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
