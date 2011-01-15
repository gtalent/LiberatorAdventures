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

var fileNotFound string = "File not found, perhaps it was taken by Tusken Raiders?"
var out *server.ChannelLine

func TopBar(signedin bool) string {
	retval, err := LoadFile("TopBar.wgt")
	if err != nil {
		return "TopBar not found."
	}
	sessionManager := "Narf!"
	if !signedin {
		if file, err := LoadFile("SessionManagerAnon.wgt"); err == nil {
			sessionManager = file
		}
	} else {
		if file, err := LoadFile("SessionManager.wgt"); err == nil {
			sessionManager = file
		}
	}
	retval = strings.Replace(retval, "{{WebHome}}", server.Settings.WebHome(), -1)
	retval = strings.Replace(retval, "{{SessionManager}}", sessionManager, -1)
	return retval
}

func postDiv() string {
	bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "post.wgt")
	if err != nil {
		return "<div><h3>{{Title}}</h3><br>{{Content}}</div>"
	}
	return string(bytes)
}

func messagePage(message string) string {
	if file, err := LoadFile("message.html"); err == nil {
		file = strings.Replace(file, "{{TopBar}}", TopBar(false), -1)
		file = strings.Replace(file, "{{Message}}", message, -1)
		return file
	}
	return fileNotFound
}

func (me *Post) HTML() string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Author}}", me.Author, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	return retval
}

func LoadFile(path string) (string, os.Error) {
	data, err := ioutil.ReadFile(server.Settings.WebRoot() + path)
	return string(data), err
}

func home(ctx *web.Context, val string) string {
	switch val {
	case "", "index.html", "index.htm":
		db, err := couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", "liberator_adventures")
		data, err := LoadFile("index.html")
		if err != nil {
			break
		}
		users := new(UserList)
		list := "<ul>\n"
		if _, err = db.Retrieve("UserList", users); err == nil {
			size := len(users.Users)
			for i := 0; i < size; i++ {
				user := users.Users[i]
				list += "<il><a href=\"" + server.Settings.WebHome() + "posts?user=" + user + "\">" + user + "</a></il><br>"
			}
		}
		list += "</ul>"
		data = strings.Replace(data, "{{TopBar}}", TopBar(false), -1)
		data = strings.Replace(data, "{{UserList}}", list, -1)
		return data
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
		_, err = db.Retrieve("BlogData_"+user, blogData)
		if err != nil {
			retval := strings.Replace(string(bytes), "{{Posts}}", "No posts from "+user+".", -1)
			retval = strings.Replace(retval, "{{TopBar}}", TopBar(false), -1)
			retval = strings.Replace(retval, "{{User}}", user, -1)
			return retval
		}
		post := new(Post)
		posts := ""
		for i := blogData.PostCount; i > 0; i-- {
			_, err := db.Retrieve("Post_"+strconv.Itoa(i)+"_"+user, post)
			if err != nil {
				post.Title = "Error: Post not found."
				post.Content = "Error: Post not found."
			} else {
				post.Content = strings.Replace(post.Content, "\n", "<br>", -1)
			}
			posts += post.HTML() + "<br>"
		}
		retval := strings.Replace(string(bytes), "{{Posts}}", posts, -1)
		retval = strings.Replace(retval, "{{TopBar}}", TopBar(false), -1)
		retval = strings.Replace(retval, "{{User.Name}}", user, -1)
		return retval

	default:
		retval, err := LoadFile(val)
		if err != nil {
			break
		}
		if strings.HasSuffix(val, ".wgt") || strings.HasSuffix(val, ".html") {
			retval = strings.Replace(retval, "{{TopBar}}", TopBar(false), -1)
		}
		return retval
	}
	return "Page not found, perhaps it was taken by Tusken Raiders?"
}

func post(ctx *web.Context, val string) string {
	switch val {
	case "CreateUser":
		username := ctx.Params["Username"]
		email := ctx.Params["Email"]
		password := ctx.Params["Password"]
		password2 := ctx.Params["Password2"]
		if password != password2 {
			return messagePage("Passwords do not match.")
		}
		user := new(User)
		user.Username = username
		user.ID = "User_" + username
		user.Email = email
		user.Password = password
		db, err := couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", "liberator_adventures")
		if err != nil {
			break
		}
		_, rev, err := db.Insert(user)
		if err != nil {
			return messagePage("Username already taken.")
		}
		users := new(UserList)
		_, err = db.Retrieve("UserList", users)
		//if can't add the user to the user list, delete the user
		if err != nil {
			db.Delete(user.ID, rev)
			return messagePage("Error")
		}

		users.Users = append(users.Users, user.Username)
		rev, err = db.Edit(users)

		//if can't add the user to the user list, delete the user
		if err != nil {
			db.Delete(user.ID, rev)
			return messagePage(err.String())
		}

		//return news of success
		if file, err := LoadFile("userCreated.html"); err == nil {
			file = strings.Replace(file, "{{TopBar}}", TopBar(false), -1)
			file = strings.Replace(file, "{{User.Name}}", username, -1)
			return file
		} else {
			break
		}
	}
	return fileNotFound
}

type dummy struct{}

func (me dummy) Write(p []byte) (n int, err os.Error) {
	return 0, nil
}

func RunWebServer(line *server.ChannelLine) {
	out = line
	var s web.Server
	s.Logger = log.New(new(dummy), "", 0)
	s.Get("/Liberator/(.*)", home)
	s.Post("/Liberator/(.*)", post)
	s.Run("0.0.0.0:" + strconv.Uitoa(server.Settings.WebPort()))
}
