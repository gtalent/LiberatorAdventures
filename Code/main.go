/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"fmt"
	"flag"
)

func main() {
	dbinit := flag.Bool("initDB", false, "Initialize the database, and then end execution.")
	settings := flag.String("conf", "blinzd.conf", "The location of the configuration file.")
	p := flag.Bool("p", false, "Indicates whether or not the program should print output to the terminal.")
	flag.Parse()

	if err := Settings.Load(*settings); err != nil {
		if *p {
			fmt.Println(err.String())
		}
		return
	}

	if *dbinit {
		fmt.Println(initDB())
		return
	}

	mainChan := make(chan string)
	webChan := NewChannelLine("Web", mainChan)
	go RunWebServer(webChan)

	if *p {
		for {
			fmt.Println(<-mainChan)
		}
	} else {
		for {
			<-mainChan
		}
	}
}
