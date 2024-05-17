package main

import (
	"embed"
	"fionna/cmd"
)

//go:embed fionna-web/dist/* fionna-web/dist/assets/*
var dist embed.FS

//go:embed fionna-web/dist/index.html
var indexHtml []byte

//func main() {
//	r := gin.Default()
//	r.Use(server.Cors())
//
//	r.GET("/", func(c *gin.Context) {
//		c.Data(http.StatusOK, "text/html", indexHtml)
//	})
//
//	r.NoRoute(func(c *gin.Context) {
//		//fmt.Println(path.Join("dist", c.Request.URL.Path))
//		data, err := dist.ReadFile(path.Join("fionna-web/dist", c.Request.URL.Path))
//
//		if err != nil {
//			//c.Redirect(http.StatusMovedPermanently, "/")
//		}
//		mimeType := mime.DetectFilePath(c.Request.URL.Path)
//		c.Data(http.StatusOK, mimeType, data)
//	})
//
//	server.InitDB()
//
//	server.GroupAndroidSerialUrl(r)
//	server.GroupAndroidPackageUrl(r)
//	server.WebSocketScrcpy(r)
//	server.WebSocketPerf(r)
//	server.WebSocketTerminal(r)
//	server.GroupReportUrl(r)
//	// 开发时可以注掉
//	gin.SetMode(gin.ReleaseMode)
//
//	port := "3417"
//	link := "http://127.0.0.1:" + port
//	fmt.Fprintf(os.Stdout, "link: \033]8;;%s\033\\%s\033]8;;\033\\\n", link, link)
//	r.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
//}

func main() {
	cmd.SetEmbed(dist, indexHtml)
	cmd.Execute()
}
