package main

import (
	"embed"
	"fionna/cmd"
	"fionna/server"
)

//go:embed fionna-web/dist/* fionna-web/dist/assets/*
var dist embed.FS

//go:embed fionna-web/dist/index.html
var indexHtml []byte

func main() {
	server.SetEmbed(dist, indexHtml)
	cmd.Execute()
}
