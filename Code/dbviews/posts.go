package main

var design_posts designDoc = designDoc{
	ID:   "_design/posts",
	Lang: "javascript",
	Views: view("by_owner", "function(doc) { if (doc.Type == 'Post')  emit(doc.Title, doc) }")}
