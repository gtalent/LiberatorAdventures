package main

import (
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
	post := NewPost()
	postID, ok := ctx.Params["PostID"]
	var newPost bool
	if ok && postID != "NewPost" {
		db.Retrieve(postID, &post)
		if userKey, ok := readUserKey(ctx); !(ok && cookies.UserKeys[userKey] == post.Owner) {
			return messagePage("You do not have permission to edit this post.", ctx)
		}
		newPost = false
	} else {
		postID = "NewPost"
		newPost = true
	}
	if file, err := LoadTemplate("", "EditPost.html", ctx); err == nil {
		if newPost {
			file = strings.Replace(file, "{{Message}}", "<h3>Writing New Post</h3>", 1)
		} else {
			file = strings.Replace(file, "{{Message}}", "<h3>Editing Existing Post</h3>", 1)
		}
		file = strings.Replace(file, "{{PostID}}", postID, 1)
		file = strings.Replace(file, "{{Title}}", post.Title, 1)
		file = strings.Replace(file, "{{Author}}", post.Author, 1)
		file = strings.Replace(file, "{{Content}}", post.Content, 1)
		authors := ""
		defaultAuthor := post.Author
		if defaultAuthor != "" {
			char := NewCharacter()
			db.Retrieve(defaultAuthor, &char)
			authors += "\t\t<option value=\"" + defaultAuthor + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}
		authors += "\t\t<option value=\"\">Me</option>\n"
		blog := NewBlogData()
		db.Retrieve("BlogData_" + readUsername(ctx), &blog)
		for i := 0; i < len(blog.Characters); i++ {
			if blog.Characters[i] != defaultAuthor {
			char := NewCharacter()
			db.Retrieve(blog.Characters[i], &char)
			authors += "\t\t<option value=\"" + blog.Characters[i] + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}}
		file = strings.Replace(file, "{{AuthorOptions}}", authors, 1)

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
	post := NewPost()
	post.ID = ctx.Params["PostID"]
	newPost := post.ID == "NewPost"
	if !newPost {
		db.Retrieve(post.ID, &post)
	}
	pleaseSignIn := "You must sign in to post."
	username := ""
	//authenticate the user
	if userkey, ok := readUserKey(ctx); !ok { //is the user signed in?
		return messagePage(pleaseSignIn, ctx)
	} else if username, ok = cookies.UserKeys[userkey]; !ok {
		return messagePage(pleaseSignIn, ctx)
	} else if post.ID != "NewPost" { //if it is not a new post, make sure the user has the right to edit it
		db.Retrieve(post.ID, &post)
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
		//manage the BlogData
		blogData := NewBlogData()
		db.Retrieve("BlogData_"+username, &blogData)
		blogData.PostIndex++
		post.ID = "Post_" + strconv.Itoa(blogData.PostIndex) + "_" + username
		blogData.Posts = append(blogData.Posts, post.ID)
		db.Edit(&blogData)
		db.Insert(&post)
	} else {
		db.Edit(&post)
	}
	return messagePage("Post saved.", ctx)
}

func viewPost(ctx *web.Context, val string) string {
	db, err := getDB()
	if err != nil {
		return fileNotFound
	}
	user := ctx.Params["user"]
	retval, err := LoadTemplate(user, "posts.html", ctx)
	if err != nil {
		return fileNotFound
	}
	blogData := NewBlogData()
	_, err = db.Retrieve("BlogData_"+user, &blogData)
	if err != nil {
		retval = strings.Replace(retval, "{{Posts}}", "No posts from "+user+".", -1)
		return retval
	}
	post := NewPost()
	posts := ""
	for i := len(blogData.Posts) - 1; i > -1; i-- {
		_, err := db.Retrieve(blogData.Posts[i], &post)
		if err != nil {
			post.Title = "Error: Post not found."
			post.Content = "Error: Post not found."
		} else {
			post.Content = strings.Replace(post.Content, "\n", "<br>", -1)
		}
		posts += post.HTML(ctx) + "<br>"
	}
	retval = strings.Replace(retval, "{{Posts}}", posts, -1)
	return retval
}

func managePostsPost(ctx *web.Context, val string) string {

	return fileNotFound
}
