package touch_test

import (
	"fionna/android/android_util"
	"fionna/android/gadb"
	touch2 "fionna/android/touch"
	"fionna/entity"
	"testing"
	"time"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestTouch_Touch(t *testing.T) {
	SetClient()
	device, err := android_util.GetDevice(client, "emulator-5554")
	if err != nil {
		panic(err)
	}
	touch := touch2.NewTouch(device)
	touch.Touch(entity.TouchInfo{
		X:         0.5,
		Y:         0.5,
		TouchType: entity.TOUCH_DOWN,
		FingerID:  0,
	})
	time.Sleep(3 * time.Second)
	touch.Touch(entity.TouchInfo{
		X:         0.5,
		Y:         0.5,
		TouchType: entity.TOUCH_UP,
		FingerID:  0,
	})
}
