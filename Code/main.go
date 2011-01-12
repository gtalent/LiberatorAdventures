package main

import (
	//"couch-go.googlecocde.com/hg"
	"blinz/server"
	"fmt"
)

func main() {
	server.Settings.Load()
	mainChan := make(chan string)
	webChan := server.NewChannelLine("Web", mainChan)
	RunWebServer(webChan)
	for {
		fmt.Println(<-mainChan)
	}
}
