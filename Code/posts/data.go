package main

//Database views

var design_posts designDoc = designDoc{ID:   "_design/posts",Lang: "javascript",
	Views: view("by_owner", "function(doc) { if (doc.Type == 'Post')  emit(doc.Title, doc) }")}

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

