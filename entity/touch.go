package entity

type TouchType string

const (
	TOUCH_DOWN = "down"
	TOUCH_MOVE = "move"
	TOUCH_UP   = "up"
)

type TouchInfo struct {
	X         float32   `json:"x"`
	Y         float32   `json:"y"`
	TouchType TouchType `json:"touchType"`
	FingerID  int       `json:"fingerID"`
}
