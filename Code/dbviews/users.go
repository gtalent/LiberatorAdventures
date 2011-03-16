package main

var design_users designDoc = designDoc{
	ID:   "_design/users",
	Lang: "javascript",
	Views: view("all", "function(doc) {if (doc.Type == 'User') emit(doc.Username, doc)}")}
