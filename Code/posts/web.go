/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package posts

import (
	"web"
	"strings"
	"strconv"
	"libadv/util"
	"libadv/char"
	"libadv/html"
)

//The HTTP get method for getting the page for editing posts.
func GetEditPost(ctx *web.Context, val string) string {
	db, err := util.GetDB()
	if err != nil {
		return util.FileNotFound
	}
	post := NewPost()
	postID, ok := ctx.Params["PostID"]
	var newPost bool
	if ok && postID != "NewPost" {
		db.Retrieve(postID, &post)
		if userKey, ok := util.ReadUserKey(ctx); !(ok && util.GetUserKey(userKey) == post.Owner) {
			return util.MessagePage("You do not have permission to edit this post.", ctx)
		}
		newPost = false
	} else {
		postID = "NewPost"
		newPost = true
	}
	if file, err := util.LoadTemplate("", "EditPost.html", ctx); err == nil {
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
			char := char.NewCharacter()
			db.Retrieve(defaultAuthor, &char)
			authors += "\t\t<option value=\"" + defaultAuthor + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}
		authors += "\t\t<option value=\"\">Me</option>\n"
		blog := util.NewBlogData()
		db.Retrieve("BlogData_" + util.ReadUsername(ctx), &blog)
		for i := 0; i < len(blog.Characters); i++ {
			if blog.Characters[i] != defaultAuthor {
			char := char.NewCharacter()
			db.Retrieve(blog.Characters[i], &char)
			authors += "\t\t<option value=\"" + blog.Characters[i] + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}}
		file = strings.Replace(file, "{{AuthorOptions}}", authors, 1)

		return file
	}
	return util.FileNotFound
}

//The HTTP post method for editing posts.
func PostEditPost(ctx *web.Context, val string) string {
	db, err := util.GetDB()
	if err != nil {
		return util.FileNotFound
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
	if userkey, ok := util.ReadUserKey(ctx); !ok { //is the user signed in?
		return util.MessagePage(pleaseSignIn, ctx)
	} else if username = util.GetUserKey(userkey); username == "" {
		return util.MessagePage(pleaseSignIn, ctx)
	} else if post.ID != "NewPost" { //if it is not a new post, make sure the user has the right to edit it
		db.Retrieve(post.ID, &post)
		if ok && post.Owner != username {
			return util.MessagePage("You do not have permission to edit this post.", ctx)
		}
	}
	//save the post
	post.Title = ctx.Params["Title"]
	post.Author = ctx.Params["Author"]
	post.Content = ctx.Params["Content"]
	post.Owner = username
	if newPost {
		//manage the BlogData
		blogData := util.NewBlogData()
		db.Retrieve("BlogData_"+username, &blogData)
		blogData.PostIndex++
		post.ID = "Post_" + strconv.Itoa(blogData.PostIndex) + "_" + username
		blogData.Posts = append(blogData.Posts, post.ID)
		db.Edit(&blogData)
		db.Insert(&post)
	} else {
		db.Edit(&post)
	}
	return util.MessagePage("Post saved.", ctx)
}

func ViewPost(ctx *web.Context, val string) string {
	db, err := util.GetDB()
	if err != nil {
		return util.FileNotFound
	}
	user := ctx.Params["user"]
	retval, err := util.LoadTemplate(user, "posts.html", ctx)
	if err != nil {
		return util.FileNotFound
	}
	blogData := util.NewBlogData()
	_, err = db.Retrieve("BlogData_"+user, &blogData)
	if err != nil {
		return util.FileNotFound
	}
	post := NewPost()
	posts := ""
	for i := len(blogData.Posts) - 1; i > -1; i-- {
		_, err := db.Retrieve(blogData.Posts[i], &post)
		if err != nil {
			post.Title = "Error: Post not found."
			post.Content = ""
		} else {
			post.Content = strings.Replace(post.Content, "\n", "<br>", -1)
		}
		posts += post.HTML(ctx) + "<br>"
	}
	retval = strings.Replace(retval, "{{Posts}}", posts, -1)
	return retval
}

func (me *Post) HTML(ctx *web.Context) string {
	retval := util.PostDiv()
	retval = strings.Replace(retval, "{{Title}}", me.Title, -1)
	if len(me.Author) != 0 {
		char := char.NewCharacter()
		db, err := util.GetDB()
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
	if username := util.ReadUsername(ctx); me.Owner == username {
		ownerControls := html.TextLink("Edit", "EditPost.html?PostID="+me.ID)
		retval = strings.Replace(retval, "{{OwnerControls}}", ownerControls.String(), -1)
	} else {
		retval = strings.Replace(retval, "{{OwnerControls}}", "", -1)
	}
	return retval
}
