package main

var design_users designDoc = designDoc{
	ID:   "_design/users",
	Lang: "javascript",
	Views: `{
      	"all": {
           "map": "function(doc) { if (doc.Type == 'User')  emit(doc.Username, doc) }"
       }
   }`}
