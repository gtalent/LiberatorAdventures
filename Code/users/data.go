package main

var design_users designDoc = designDoc{ID:   "_design/users", Lang: "javascript",
	Views: view("all", "function(doc) {if (doc.Type == 'User') emit(doc.Username, doc)}")}

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
