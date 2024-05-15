package entity

import "time"

type SerialInfo struct {
	UUID          string    `json:"uuid" gorm:"primaryKey"`
	SerialName    string    `json:"udid" gorm:"udid"`
	Name          string    `json:"name" gorm:"-"`
	TestName      *string   `json:"testName,omitempty" gorm:"testName"`
	Size          string    `json:"size" gorm:"size"`
	CPUArch       string    `json:"cpu" gorm:"cpu"`
	Voltage       int       `json:"voltage" gorm:"-"`
	ProductDevice string    `json:"productDevice" gorm:"productDevice"`
	Model         string    `json:"model" gorm:"model"`
	Temperature   int       `json:"temperature" gorm:"-"`
	Manufacturer  string    `json:"manufacturer" gorm:"manufacturer"`
	Platform      int       `json:"platform"  gorm:"platform"` // default 1
	Version       string    `json:"version" gorm:"version"`
	IsHm          int       `json:"isHm"  gorm:"isHm"`
	Level         int       `json:"level" gorm:"-"` // 电池电量
	Timestamp     *int64    `json:"timestamp,omitempty" `
	PackageName   *string   `json:"packageName,omitempty"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

//type BaseModel struct {
//	UUID      string `gorm:"primaryKey"`
//	TestName  string `gorm:"testName"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
