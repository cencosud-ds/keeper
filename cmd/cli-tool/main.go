package main

import (
	"keeper/pkg/cli-tool"
	"log"
)

func main() {
	// Makes error line shows on "log" lib usage
	log.SetFlags(log.LstdFlags | log.Llongfile)

	cli_tool.Execute()
}
