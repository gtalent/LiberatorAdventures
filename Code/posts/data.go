package main

import (
	"strings"
	"web"
	"blinz/html"
)

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

func (me *Post) HTML(ctx *web.Context) string {
	retval := postDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	if len(me.Author) != 0 {
		char := NewCharacter()
		db, err := getDB()
		if err == nil {
			db.Retrieve(me.Author, &char)
			retval = strings.Replace(retval, "{{Author}}", "<a href=\"Character.html?CharID=" + me.Author + "\">" + char.Name + "</a>", -1)
		} else {
			retval = strings.Replace(retval, "{{Author}}", "", -1)
		}
	} else {
		retval = strings.Replace(retval, "{{Author}}", me.Owner, -1)
	}
	retval = strings.Replace(retval, "{{Content}}", me.Content, -1)
	if username := readUsername(ctx); me.Owner == username {
		ownerControls := html.TextLink("Edit", "EditPost.html?PostID="+me.ID)
		retval = strings.Replace(retval, "{{OwnerControls}}", ownerControls.String(), -1)
	} else {
		retval = strings.Replace(retval, "{{OwnerControls}}", "", -1)
	}
	return retval
}
