package scrcpy_client

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"net"
)

type Control struct {
	controlConn net.Conn
}

func NewControl(controlConn net.Conn) *Control {
	return &Control{controlConn: controlConn}
}

//func (c *Control) KeyCode(keyCode, action, repeat int) error {
//	var data = []interface{}{
//		uint8(TypeInjectKeycode), // base
//		uint8(action),            // B 1 byte
//		uint32(keyCode),          // i 4 byte
//		uint32(repeat),           // i 4 byte
//		uint32(0),                // i 4 byte
//	}
//	msg, err := serializePack(data)
//	if err != nil {
//		return err
//	}
//	_, err = c.controlConn.Write(msg)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (c *Control) Text(text string) error {
	var data = []interface{}{
		uint8(TypeInjectTEXT),     // base
		uint32(len([]byte(text))), // i 4 byte
	}
	msg, err := serializePack(data)
	if err != nil {
		log.Error("control serialize text packet err", err)
		return err
	}
	msg = append(msg, []byte(text)...)
	_, err = c.controlConn.Write(msg)
	if err != nil {
		log.Error("control send text packet err", err)
		return err
	}
	return nil
}

//func (c *Control) Touch(touch *entity.ScrcpyTouch, touchID int) error {
//	var data = []interface{}{
//		uint8(TypeInjectTOUCHEvent), // base       1
//		uint8(touch.ActionType),     // B 1 byte   2
//		int64(touchID),              // q 8 byte   10
//		uint32(touch.X),             // i 4 byte   14
//		uint32(touch.Y),             // i 4 byte   18
//		uint16(touch.Width),         // H 2 byte   20
//		uint16(touch.Height),        // H 2 byte   22
//		uint16(0xffff),              // H 2 byte   24
//		uint32(1),                   // i 4 byte   28
//		uint32(1),                   // i 4 byte   32
//	}
//	msg, err := serializePack(data)
//	if err != nil {
//		log.Error("control serialize touch packet err", err)
//		return err
//	}
//	_, err = c.controlConn.Write(msg)
//	if err != nil {
//		log.Error("control send touch packet err", err)
//		return err
//	}
//	return nil
//}

func serializePack(data []interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

const (
	TypeInjectKeycode           = 0
	TypeInjectTEXT              = 1
	TypeInjectTOUCHEvent        = 2
	TypeInjectSCROLLEvent       = 3
	TypeBACKORScreenON          = 4
	TypeEXPANDNOTIFICATIONPANEL = 5
	TypeEXPANDSETTINGSPANEL     = 6
	TypeCOLLAPSEPANELS          = 7
	TypeGETCLIPBOARD            = 8
	TypeSETCLIPBOARD            = 9
	TypeSETScreenPowerMode      = 10
	TypeROTATEDEVICE            = 11
)
