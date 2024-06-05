package cmd

import (
	"embed"
	"fionna/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qingstor/go-mime"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path"
)

func runWeb() {
	r := gin.Default()
	r.Use(server.Cors())

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

	server.InitDB(dbName)

	server.GroupAndroidSerialUrl(r)
	server.GroupAndroidPackageUrl(r)
	server.WebSocketScrcpy(r)
	server.WebSocketPerf(r)
	server.WebSocketTerminal(r)
	server.GroupReportUrl(r)
	// 开发时可以注掉
	gin.SetMode(gin.ReleaseMode)

	port := "3417"
	link := fmt.Sprintf("http://127.0.0.1:%s", port)
	fmt.Fprintf(os.Stdout, "link: \033]8;;%s\033\\%s\033]8;;\033\\\n", link, link)
	r.Run(fmt.Sprintf(":%s", port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var (
	dist      embed.FS
	indexHtml []byte
	dbName    string
)

func SetEmbed(distParam embed.FS, indexHtmlParam []byte) {
	dist = distParam
	indexHtml = indexHtmlParam
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Fionna web mode",
	Long:  "Fionna web mode",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		runWeb()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&dbName, "db-name", "d", "test.db", "specify the sql lite name to use")
}
