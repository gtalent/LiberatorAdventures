package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"os"
	"web"
	"couch-go.googlecode.com/hg"
	"blinz/server"
)

var out *server.ChannelLine
var TopBar string
var postDiv string

type BlogData struct {
	PostCount int
}

type Post struct {
	Title, Author, Date, Content string
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
}

func home(ctx *web.Context, val string) string {
	switch val {
	case "posts", "posts/":
		db, err := couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", "liberator_adventures")
		if err != nil {
			break
		}
		bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "posts.html")
		if err != nil {
			break
		}
		blogData := new(BlogData)
		user := ctx.Params["user"]
		_, err = db.Retrieve("BlogData_" + user, blogData)
		if err != nil {
			retval := strings.Replace(string(bytes), "{{Posts}}", "No posts from " + user + ".", -1)
			retval = strings.Replace(retval, "{{TopBar}}", TopBar, -1)
			return retval
		}
		post := new(Post)
		posts := ""
		for i := blogData.PostCount; i > 0; i-- {
			_, err := db.Retrieve("Post_" + strconv.Itoa(i) + "_" + user, post)
			if err != nil {
				post.Title = "Error: Post not found."
				post.Content = "Error: Post not found."
			} else {
				post.Content = strings.Replace(post.Content, "\n", "<br>", -1)
			}
			posts += post.HTML() + "<br>"
		}
		retval := strings.Replace(string(bytes), "{{Posts}}", posts, -1)
		retval = strings.Replace(retval, "{{TopBar}}", TopBar, -1)
		return retval
	}
	return "Page not found, perhaps it was taken by Tusken Raiders?"
}

type dummy struct{}

func (me dummy) Write(p []byte) (n int, err os.Error) {
	return 0, nil
}

func RunWebServer(line *server.ChannelLine) {
	out = line
	load()
	var s web.Server
	s.Logger = log.New(new(dummy), "", 0)
	s.Get("/liberator-blog/(.*)", home)
	s.Run("0.0.0.0:" + strconv.Uitoa(server.Settings.WebPort()))
}
