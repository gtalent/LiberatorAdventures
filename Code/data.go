package main

import (
	"strings"
	"os"
	"web"
	"blinz/html"
	"blinz/server"
	"couch-go.googlecode.com/hg"
)

var cookies *Cookies = NewCookies()

//Gets the database connection.
func getDB() (couch.Database, os.Error) {
	return couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", server.Settings.Database())
}

type User struct {
	ID                        string "_id"
	Rev                       string "_rev"
	Type                      string
	Username, Email, Password string
}

//Returns a new User object by value.
func NewUser() User {
	var data User
	data.Type = "User"
	return data
}

type BlogData struct {
	ID             string "_id"
	Rev            string "_rev"
	Type           string
	CharacterIndex int
	Characters     []string
	PostIndex      int "PostCount"
	Posts          []string
}

//Returns a new BlogData object by value.
func NewBlogData() BlogData {
	var data BlogData
	data.Type = "BlogData"
	return data
}

type Cookies struct {
	ID       string "_id"
	Rev      string "_rev"
	Type     string
	UserKeys map[string]string
}

func NewCookies() *Cookies {
	retval := new(Cookies)
	retval.UserKeys = make(map[string]string)
	return retval
}

type Post struct {
	ID                                  string "_id"
	Rev                                 string "_rev"
	Type                                string
	Title, Author, Owner, Date, Content string
}

//Returns a new Post object by value.
func NewPost() Post {
	var data Post
	data.Type = "Post"
	return data
}

func (me *Post) HTML(ctx *web.Context) string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Author}}", me.Author, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	if username := readUsername(ctx); me.Owner == username {
		ownerControls := html.TextLink("Edit", "EditPost.html?PostID="+me.ID)
		retval = strings.Replace(retval, "{{OwnerControls}}", ownerControls.String(), -1)
	} else {
		retval = strings.Replace(retval, "{{OwnerControls}}", "", -1)
	}
	return retval
}


type Character struct {
	ID                                        string "_id"
	Rev                                       string "_rev"
	Type                                      string
	Game, Name, World, Alligiance, Bio, Owner string
}

func NewCharacter() Character {
	var data Character
	data.Type = "Post"
	return data

}

//Unfinished
//Checks the username to make sure it does not contain any illegal characters.
//Returns a message if the name is bad, and a boolean indicating whether or not the name was ok
func checkUsername(username string) (string, bool) {
	//badChars := ";\\ \"\n\t={}()[]"
	if strings.Contains(username, ";") || strings.Contains(username, "\\") || strings.Contains(username, " ") || strings.Contains(username, "=") {
		return "Usernames may not contian the following characters: (A-Z, a-z, 0-9, _, -)", false
	}
	return "", true
}
