package main

import (
	"embed"
	"fionna/cmd"
)

//go:embed fionna-web/dist/* fionna-web/dist/assets/*
var dist embed.FS

//go:embed fionna-web/dist/index.html
var indexHtml []byte

func main() {
	cmd.SetEmbed(dist, indexHtml)
	cmd.Execute()
}
