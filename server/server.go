package server

import (
	"fionna/android/gadb"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os/exec"
)

var (
	db     *gorm.DB
	client gadb.Client
)

var isSQLDebug = false

func SetSQLDebug(flag bool) {
	isSQLDebug = flag
}

func init() {
	var err error
	client, err = gadb.NewClient()
	if err != nil {
		cmd := exec.Command("adb", "start-server")
		cmd.Run()
		var err1 error
		client, err1 = gadb.NewClient()
		if err1 != nil {
			panic("failed to connect adb server")
		}
	}
}

func InitDB() {
	var err error
	var gConfig = &gorm.Config{}

	if !isSQLDebug {
		gConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err = gorm.Open(sqlite.Open("test.db"), gConfig)
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

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// solve cross domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//func Server() {
//	InitDB()
//	r := gin.Default()
//	r.Use(Cors())
//
//	//pprof.Register(r)
//	GroupAndroidSerialUrl(r)
//	GroupAndroidPackageUrl(r)
//	//GroupScrcpy(r)
//	WebSocketScrcpy(r)
//	WebSocketPerf(r)
//	WebSocketTerminal(r)
//	GroupReportUrl(r)
//	//WebSocketScrcpy1(r)
//
//	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
//}
