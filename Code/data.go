package main

import (
	"strings"
	"web"
	"blinz/html"
)

type User struct {
	ID                        string "_id"
	Rev                       string "_rev"
	Username, Email, Password string
}

type UserList struct {
	ID    string "_id"
	Rev   string "_rev"
	Users []string
}

type BlogData struct {
	ID        string "_id"
	Rev       string "_rev"
	PostIndex int "PostCount"
	Posts	  []string
}

type Post struct {
	ID                                  string "_id"
	Rev                                 string "_rev"
	Title, Author, Owner, Date, Content string
}

func (me *Post) HTML(ctx *web.Context) string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Author}}", me.Author, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	if username := readUsername(ctx); me.Owner == username {
		ownerControls := html.TextLink("Edit", "EditPost.html?PostID=" + me.ID)
		retval = strings.Replace(retval, "{{OwnerControls}}", ownerControls.String(), -1)
	} else {
		retval = strings.Replace(retval, "{{OwnerControls}}", "", -1)
	}
	return retval
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
