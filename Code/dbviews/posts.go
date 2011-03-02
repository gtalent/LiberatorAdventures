package main

var design_posts designDoc = designDoc{
	ID:   "_design/posts",
	Lang: "javascript",
	Views: `{
       "by_owner": {
           "map": "function(doc) { if (doc.Type == 'Post')  emit(doc.Title, doc) }"
       }
   }`}
