package main

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
	PostCount int
}

type Post struct {
	Title, Author, Date, Content string
}
