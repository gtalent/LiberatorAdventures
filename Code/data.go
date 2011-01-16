package main

import "strings"

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
	PostCount int
}

type Post struct {
	ID                                  string "_id"
	Rev                                 string "_rev"
	Title, Author, Owner, Date, Content string
}

func (me *Post) HTML() string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Author}}", me.Author, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
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
