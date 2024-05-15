package entity

type TerminalRecvMessageType string
type TerminalSendMessageType string

const (
	CommandTerminalType    TerminalRecvMessageType = "command"
	LogcatTerminalType     TerminalRecvMessageType = "logcat"
	StopCommandType        TerminalRecvMessageType = "stopCommand"
	StopLogcatType         TerminalRecvMessageType = "stopLogcat"
	PongTerminalType       TerminalRecvMessageType = "pongTerminal"
	CloseTerminalType      TerminalRecvMessageType = "closeTerminal"
	SendCommandRespType    TerminalSendMessageType = "commandResp"
	SendCommandRespEndType TerminalSendMessageType = "commandRespEnd"
	SendLogcatRespType     TerminalSendMessageType = "logcatResp"
	SendLogcatRespEndType  TerminalSendMessageType = "logcatRespEnd"
	TerminalError          TerminalSendMessageType = "error"
)

type TerminalRecvMessage struct {
	MessageType TerminalRecvMessageType `json:"messageType"`
	Uuid        string                  `json:"uuid"`
	Data        interface{}             `json:"data"`
}

func NewTerminalErrorInfo(message string) *TerminalErrorInfo {
	return &TerminalErrorInfo{
		MessageType: TerminalError,
		ErrorInfo:   message,
	}
}

type TerminalErrorInfo struct {
	MessageType TerminalSendMessageType `json:"messageType"`
	ErrorInfo   string                  `json:"errorInfo"`
}

type TerminalSendMessage struct {
	MessageType TerminalSendMessageType `json:"messageType"`
	Uuid        string                  `json:"uuid"`
	Data        interface{}             `json:"data"`
}

type LogcatRecvMessage struct {
	Level  string `json:"level"`
	Filter string `json:"filter"`
}
