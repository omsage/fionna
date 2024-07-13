package cmd

import (
	"embed"
	"fionna/server"
	"github.com/spf13/cobra"
)

func runWeb() {
	server.Init(dbName)
}

var (
	dist      embed.FS
	indexHtml []byte
	dbName    string
)

//func SetEmbed(distParam embed.FS, indexHtmlParam []byte) {
//	dist = distParam
//	indexHtml = indexHtmlParam
//}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Fionna web mode",
	Long:  "Fionna web mode",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		server.Init(dbName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&dbName, "db-path", "d", "test.db", "specify the SQLite path to use")
}
