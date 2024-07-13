package entity

type TouchType string

const (
	TOUCH_DOWN = "down"
	TOUCH_MOVE = "move"
	TOUCH_UP   = "up"
)

type TouchInfo struct {
	X         float32
	Y         float32
	TouchType TouchType
	FingerID  int
}
