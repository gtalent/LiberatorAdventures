package main

import "strings"

type User struct {
	ID string "_id"
	Rev string "_rev"
	Username, Email, Password string
}

type UserList struct {
	ID string "_id"
	Rev string "_rev"
	Users []string
}

type BlogData struct {
	ID string "_id"
	Rev string "_rev"
	PostCount int
}

type Post struct {
	Title, Author, Date, Content string
}

func (me *Post) HTML() string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	retval = strings.Replace(retval, "{{Author}}", me.Author, -1)
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	return retval
}
