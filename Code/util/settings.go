/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
//File for loading and managing the server configuration.
package main

import (
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

var Settings ServerConf

func parseSetting(line, setting string) string {
	return strings.Trim(strings.Replace(line, setting + ":", "", -1), " \t")
}

//Conveniently loads and holds all settings information.
type ServerConf struct {
	webHome, webRoot, databaseAddress, database, cookieSecret string
	webPort         uint
}

//The cookie secret used to ensure secure transmission of cookies.
func (me *ServerConf) CookieSecret() string {
	return me.cookieSecret
}

//The address at which to find the database.
func (me *ServerConf) DatabaseAddress() string {
	return me.databaseAddress
}

//The name of the database to connect to.
func (me *ServerConf) Database() string {
	return me.database
}

//The root directory to serve the website from.
func (me *ServerConf) WebRoot() string {
	return me.webRoot
}

//The root directory of the website from the perspective of the client.
func (me *ServerConf) WebHome() string {
	return me.webHome
}

//The port that this runs as an HTTP server.
func (me *ServerConf) WebPort() uint {
	return me.webPort
}

func (me *ServerConf) Load(path string) os.Error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	list := strings.Split(string(bytes), "\n", -1)
	for i := 0; i < len(list); i++ {
		if strings.HasPrefix(list[i], "DatabaseAddress:") {
			me.databaseAddress = parseSetting(list[i], "DatabaseAddress")
		} else if strings.HasPrefix(list[i], "CookieSecret:") {
			me.cookieSecret = parseSetting(list[i], "CookieSecret")
		} else if strings.HasPrefix(list[i], "Database:") {
			me.database = parseSetting(list[i], "Database")
		} else if strings.HasPrefix(list[i], "WebHome:") {
			me.webHome = parseSetting(list[i], "WebHome")
		} else if strings.HasPrefix(list[i], "WebRoot:") {
			me.webRoot = parseSetting(list[i], "WebRoot")
		} else if strings.HasPrefix(list[i], "WebPort:") {
			me.webPort, err = strconv.Atoui(parseSetting(list[i], "WebPort"))
			if err != nil {
				me.webPort = 8080
			}
		}
	}
	return nil
}
