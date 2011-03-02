package main

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
