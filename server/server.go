package server

import (
	"embed"
	"fionna/server/android"
	"fionna/server/db"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qingstor/go-mime"
	"net/http"
	"os"
	"path"
	"runtime"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// solve cross domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var (
	dist      embed.FS
	indexHtml []byte
)

func SetEmbed(distParam embed.FS, indexHtmlParam []byte) {
	dist = distParam
	indexHtml = indexHtmlParam
}

func Init(dbName string) {
	db.InitDB(dbName)
	android.Init(upGrader)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", indexHtml)
	})

	r.NoRoute(func(c *gin.Context) {
		//fmt.Println(path.Join("dist", c.Request.URL.Path))
		data, err := dist.ReadFile(path.Join("fionna-web/dist", c.Request.URL.Path))

		if err != nil {
			//c.Redirect(http.StatusMovedPermanently, "/")
		}
		mimeType := mime.DetectFilePath(c.Request.URL.Path)
		c.Data(http.StatusOK, mimeType, data)
	})
	// 初始化安卓服务
	AndroidServerInit(r)
	GroupReportUrl(r)
	// 开发时可以注掉
	gin.SetMode(gin.ReleaseMode)

	port := "3417"
	link := fmt.Sprintf("http://127.0.0.1:%s", port)
	if runtime.GOOS == "windows" {
		// Windows 下的处理方式
		fmt.Fprintf(os.Stdout, "link: %s\n", link)
	} else {
		// 其他平台下的 ANSI Escape Code 处理
		fmt.Fprintf(os.Stdout, "link: \033]8;;%s\033\\%s\033]8;;\033\\\n", link, link)
	}
	r.Run(fmt.Sprintf(":%s", port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func AndroidServerInit(r *gin.Engine) {
	android.GroupAndroidSerialUrl(r)
	android.GroupAndroidPackageUrl(r)
	android.WebSocketScrcpy(r)
	//android.WebSocketPerf(r)
	android.WebSocketTerminal(r)
	android.AndroidControl(r)
}
