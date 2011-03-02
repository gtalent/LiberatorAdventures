package main

var design_characters designDoc = designDoc{
	ID:   "_design/characters",
	Lang: "javascript",
	Views: `{
       "all": {
           "map": "function(doc) { if (doc.Type == 'Character')  emit(doc.Name, doc) }"
       }
   }`}
