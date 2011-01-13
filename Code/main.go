package main

import (
	"blinz/server"
	//"fmt"
)

func main() {
	if err := server.Settings.Load("/usr/local/etc/LiberatorAdventuresd.conf"); err != nil {
		//fmt.Println(err.String())
		return
	}
	mainChan := make(chan string)
	webChan := server.NewChannelLine("Web", mainChan)
	go RunWebServer(webChan)
	for {
		//fmt.Println(<-mainChan)
		_ = <-mainChan
	}
}
