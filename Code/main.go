/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"fmt"
	"flag"
	"libadv/char"
	"libadv/posts"
	"libadv/util"
	"libadv/users"
)

//Initializes the database by adding the design documents.
func initDB() bool {
	db, err := util.GetDB()
	if err != nil {
		return false
	}
	_, _, err1 := db.Insert(users.Design_users)
	_, _, err2 := db.Insert(posts.Design_posts)
	_, _, err3 := db.Insert(char.Design_characters)
	return err1 == nil && err2 == nil && err3 == nil
}

func main() {
	dbinit := flag.Bool("initDB", false, "Initialize the database, and then end execution.")
	settings := flag.String("conf", "blinzd.conf", "The location of the configuration file.")
	p := flag.Bool("p", false, "Indicates whether or not the program should print output to the terminal.")
	flag.Parse()

	if err := util.Settings.Load(*settings); err != nil {
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
	webChan := util.NewChannelLine("Web", mainChan)
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
