package entity

type ServerCode int

const (
	RequestSucceed     ServerCode = 10000
	GetPerfErr         ServerCode = 10001
	ScrcpyServerErr    ServerCode = 10002
	TouchServerErr     ServerCode = 10003
	GetDeviceErr       ServerCode = 10004
	ParameterErr       ServerCode = 10005
	ServiceErr         ServerCode = 10006
	ScrcpyMessageErr   ServerCode = 10007
	TerminalMessageErr ServerCode = 10008
	PerfMessageErr     ServerCode = 10009
	CustomErr          ServerCode = 100032
)

var CodeDefaultMessage = map[ServerCode]string{
	RequestSucceed:     "Succeed",
	GetPerfErr:         "Get Perf Error",
	ScrcpyServerErr:    "Scrcpy Server Error",
	TouchServerErr:     "Touch Server Error",
	GetDeviceErr:       "Get Device Error",
	ParameterErr:       "Parameter Error",
	ServiceErr:         "Service Error",
	ScrcpyMessageErr:   "Scrcpy Message Error",
	TerminalMessageErr: "Terminal Message Error",
	PerfMessageErr:     "Perf Message Error",
}

type ResponseData struct {
	Data interface{} `json:"data"`
	Code ServerCode  `json:"code"`
}
