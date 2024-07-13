package entity

type ScrcpyRecvMessageType string
type ScrcpySendMessageType string

const (
	ScrcpyTouchType    ScrcpyRecvMessageType = "touch"
	ScrcpyCloseType    ScrcpyRecvMessageType = "close"
	ScrcpyPongType     ScrcpyRecvMessageType = "pong"
	ScrcpyDeviceType   ScrcpyRecvMessageType = "device"
	ScrcpySizeInfoType ScrcpySendMessageType = "sizeInfo"
	ScrcpyPicType      ScrcpyRecvMessageType = "pic"
	ScrcpyErrorType    ScrcpySendMessageType = "error"
)

type ScrcpyPic string

const (
	ScrcpyPicLow    ScrcpyPic = "low"
	ScrcpyPicMid    ScrcpyPic = "mid"
	ScrcpyPicHeight ScrcpyPic = "height"
)

type ScrcpyDevice struct {
	UDID string `json:"udid"`
}

//type ScrcpyTouch struct {
//	ActionType ActionType `json:"actionType"`
//	X          int        `json:"x"`
//	Y          int        `json:"y"`
//	Width      int        `json:"width"`
//	Height     int        `json:"height"`
//}
//
//type ActionType int
//
//const (
//	ActionDown ActionType = 0
//	ActionUp   ActionType = 1
//	ActionMove ActionType = 2
//)

type ScrcpySizeInfo struct {
	Rotation int `json:"rotation"`
	Width    int `json:"width"`
	Height   int `json:"height"`
}

type ScrcpyRecvMessage struct {
	MessageType ScrcpyRecvMessageType `json:"messageType"`
	Data        interface{}           `json:"data"`
}

func NewScrcpyError(message string) *ScrcpyErrorMessage {
	return &ScrcpyErrorMessage{
		MessageType: ScrcpyErrorType,
		ErrorInfo:   message,
	}
}

type ScrcpyErrorMessage struct {
	MessageType ScrcpySendMessageType `json:"messageType"`
	ErrorInfo   string                `json:"errorInfo"`
}

type ScrcpySizeInfoMessage struct {
	MessageType ScrcpySendMessageType `json:"messageType"`
	Data        *ScrcpySizeInfo       `json:"sizeInfo"`
}
