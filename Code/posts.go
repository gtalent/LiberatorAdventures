package main

import (
	"blinz/server"
	"io/ioutil"
	"web"
	"strings"
	"strconv"
)

//The HTTP get method for getting the page for editing posts.
func getEditPost(ctx *web.Context, val string) string {
	db, err := getDB()
	if err != nil {
		return fileNotFound
	}
	post := new(Post)
	postID, ok := ctx.Params["PostID"]
	var newPost bool
	if ok && postID != "NewPost" {
		db.Retrieve(postID, post)
		if userKey, ok := readUserKey(ctx); !(ok && cookies[userKey] == post.Owner) {
			return messagePage("You do not have permission to edit this post.", ctx)
		}
		newPost = false
	} else {
		postID = "NewPost"
		newPost = true
	}
	if file, err := LoadFile("EditPost.html"); err == nil {
		if newPost {
			file = strings.Replace(file, "{{Message}}", "<h3>Writing New Post</h3>", 1)
		} else {
			file = strings.Replace(file, "{{Message}}", "<h3>Editing Existing Post</h3>", 1)
		}
		file = strings.Replace(file, "{{PostID}}", postID, 1)
		file = strings.Replace(file, "{{TopBar}}", TopBar(ctx), 1)
		file = strings.Replace(file, "{{Title}}", post.Title, 1)
		file = strings.Replace(file, "{{Author}}", post.Author, 1)
		file = strings.Replace(file, "{{Content}}", post.Content, 1)
		return file
	}
	return fileNotFound
}

//The HTTP post method for editing posts.
func postEditPost(ctx *web.Context, val string) string {
	db, err := getDB()
	if err != nil {
		return fileNotFound
	}
	post := new(Post)
	post.ID = ctx.Params["PostID"]
	newPost := post.ID == "NewPost"
	if !newPost {
		db.Retrieve(post.ID, post)
	}
	pleaseSignIn := "You must sign in to post."
	username := ""
	//authenticate the user
	if userkey, ok := readUserKey(ctx); !ok { //is the user signed in?
		return messagePage(pleaseSignIn, ctx)
	} else if username, ok = cookies[userkey]; !ok {
		return messagePage(pleaseSignIn, ctx)
	} else if post.ID != "NewPost" { //if it is not a new post, make sure the user has the right to edit it
		db.Retrieve(post.ID, post)
		if ok && post.Owner != username {
			return messagePage("You do not have permission to edit this post.", ctx)
		}
	}
	//save the post
	post.Title = ctx.Params["Title"]
	post.Author = ctx.Params["Author"]
	post.Content = ctx.Params["Content"]
	post.Owner = username
	if newPost {
		blogData := new(BlogData)
		db.Retrieve("BlogData_"+username, blogData)
		blogData.PostCount++
		db.Edit(blogData)
		post.ID = "Post_" + strconv.Itoa(blogData.PostCount) + "_" + username
		db.Insert(post)
	} else {
		db.Edit(post)
	}
	return messagePage("Post saved.", ctx)
}

func viewPost(ctx *web.Context, val string) string {
	db, err := getDB()
	if err != nil {
		return fileNotFound
	}
	bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "posts.html")
	if err != nil {
		return fileNotFound
	}
	blogData := new(BlogData)
	user := ctx.Params["user"]
	_, err = db.Retrieve("BlogData_"+user, blogData)
	if err != nil {
		retval := strings.Replace(string(bytes), "{{Posts}}", "No posts from "+user+".", -1)
		retval = strings.Replace(retval, "{{TopBar}}", TopBar(ctx), -1)
		retval = strings.Replace(retval, "{{User}}", user, -1)
		return retval
	}
	post := new(Post)
	posts := ""
	for i := blogData.PostCount; i > 0; i-- {
		_, err := db.Retrieve("Post_"+strconv.Itoa(i)+"_"+user, post)
		if err != nil {
			post.Title = "Error: Post not found."
			post.Content = "Error: Post not found."
		} else {
			post.Content = strings.Replace(post.Content, "\n", "<br>", -1)
		}
		posts += post.HTML(ctx) + "<br>"
	}
	retval := strings.Replace(string(bytes), "{{Posts}}", posts, -1)
	retval = strings.Replace(retval, "{{TopBar}}", TopBar(ctx), -1)
	retval = strings.Replace(retval, "{{User.Name}}", user, -1)
	return retval
}

func managePostsPost(ctx *web.Context, val string) string {

	return fileNotFound
}