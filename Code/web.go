package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"web"
	"couch-go.googlecode.com/hg"
	"blinz/server"
)

var out *server.ChannelLine
var db couch.Database
var TopBar string
var postDiv string

type BlogData struct {
	PostCount int
}

type Post struct {
	Title, Date, Content string
}

func (me *Post) HTML() string {
	retval := postDiv
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	return retval
}

func load() {
	{
		bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "topbar.html")
		if err == nil {
			retval := string(bytes)
			retval = strings.Replace(retval, "{{WebHome}}", server.Settings.WebHome(), -1)
			TopBar = retval
		} else {
			out.Put(err.String())
		}
	}
	{
		bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "post.wgt")
		if err == nil {
			retval := string(bytes)
			retval = strings.Replace(retval, "{{WebHome}}", server.Settings.WebHome(), -1)
			postDiv = retval
		} else {
			out.Put(err.String())
		}
	}
	temp, err := couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", "liberator_adventures")
	if err != nil {
		out.Put(err.String())
	} else {
		db = temp
	}
}

func home(val string) string {
	switch val {
	case "posts", "posts/":
		bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "posts.html")
		if err != nil {
			return "Page not found, perhaps it was taken by Tusken Raiders?"
		}
		retval := strings.Replace(string(bytes), "{{Posts}}", val, -1)
		retval = strings.Replace(retval, "{{TopBar}}", TopBar, -1)
		return retval
	}
	return "Page not found, perhaps it was taken by Tusken Raiders?"
}

func RunWebServer(line *server.ChannelLine) {
	out = line
	load()
	web.Get("/(.*)", home)
	web.Run("0.0.0.0:" + strconv.Uitoa(server.Settings.WebPort()))
}
