package db

import (
	"fionna/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func GetDB() *gorm.DB {
	return db
}

var isSQLDebug = false

func SetSQLDebug(flag bool) {
	isSQLDebug = flag
}

func InitDB(dbName string) {
	var err error
	var gConfig = &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if !isSQLDebug {
		gConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err = gorm.Open(sqlite.Open(dbName), gConfig)
	if err != nil {
		panic("failed to connect database")
	}
	//
	//if !db.Migrator().HasTable(&entity.BaseModel{}) {
	//	err = db.AutoMigrate(&entity.BaseModel{})
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	if !db.Migrator().HasTable(&entity.SerialInfo{}) {
		err = db.AutoMigrate(&entity.SerialInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.PerfConfig{}) {
		err = db.AutoMigrate(&entity.PerfConfig{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.ProcCpuInfo{}) {
		err = db.AutoMigrate(&entity.ProcCpuInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.ProcMemInfo{}) {
		err = db.AutoMigrate(&entity.ProcMemInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.ProcThreadsInfo{}) {
		err = db.AutoMigrate(&entity.ProcThreadsInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemCPUInfo{}) {
		err = db.AutoMigrate(&entity.SystemCPUInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SysFrameInfo{}) {
		err = db.AutoMigrate(&entity.SysFrameInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemMemInfo{}) {
		err = db.AutoMigrate(&entity.SystemMemInfo{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemNetworkInfo{}) {
		err = db.AutoMigrate(&entity.SystemNetworkInfo{})
		if err != nil {
			panic(err)
		}
	}

	// overview
	if !db.Migrator().HasTable(&entity.SystemNetworkSummary{}) {
		err = db.AutoMigrate(&entity.SystemNetworkSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemCPUSummary{}) {
		err = db.AutoMigrate(&entity.SystemCPUSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemMemSummary{}) {
		err = db.AutoMigrate(&entity.SystemMemSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.FrameSummary{}) {
		err = db.AutoMigrate(&entity.FrameSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.ProcCpuSummary{}) {
		err = db.AutoMigrate(&entity.ProcCpuSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.ProcMemSummary{}) {
		err = db.AutoMigrate(&entity.ProcMemSummary{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SysTemperature{}) {
		err = db.AutoMigrate(&entity.SysTemperature{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&entity.SystemTemperatureSummary{}) {
		err = db.AutoMigrate(&entity.SystemTemperatureSummary{})
		if err != nil {
			panic(err)
		}
	}
}
