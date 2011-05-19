/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package	users

import (
	"web"
	"rand"
	"strings"
	"strconv"
	"libadv/char"
	"libadv/posts"
	"libadv/util"
)

//session management

func SigninPost(ctx *web.Context, val string) string {
	username := ctx.Params["Username"]
	password := ctx.Params["Password"]
	user := NewUser()
	if db, err := util.GetDB(); err == nil {
		if _, err = db.Retrieve("User_"+username, &user); err == nil {
			if string(util.PasswordHash(password, user.Password.Version).Hash) == string(user.Password.Hash) {
				num := rand.Int63()
				key := username + "_" + strconv.Itoa64(num)
				util.SetUserKey(key, username)
				ctx.SetSecureCookie("UserKey", key, 6000000)
				return util.MessagePage("You are now signed in.", ctx)
			}
			return util.MessagePage("Invalid username and password combination.", ctx)
		}
		return util.MessagePage("Error: Username may not exist.", ctx)
	}
	return util.MessagePage("Could not access the database.", ctx)
}

//account management

func AccountManagementGet(ctx *web.Context, val string) string {
	file, err := util.LoadTemplate("Accout Management", val, ctx)
	if err != nil {
		return util.FileNotFound
	}
	return file
}

func CreateAccountPost(ctx *web.Context, val string) string {
	username := ctx.Params["Username"]
	email := ctx.Params["Email"]
	password := ctx.Params["Password"]
	password2 := ctx.Params["Password2"]
	if password != password2 {
		return util.MessagePage("Passwords do not match.", ctx)
	}
	if len(password) < 6 {
		return util.MessagePage("You're password must be at least 6 characters long.", ctx)
	}
	if strings.Contains(username, ";") || strings.Contains(username, "\\") || strings.Contains(username, " ") || strings.Contains(username, "=") {
		return util.MessagePage("Usernames may not contian the following characters: \" \", \"=\", \"\\\", or \";\".", ctx)
	}
	user := NewUser()
	user.Username = username
	user.ID = "User_" + username
	user.Email = email
	user.Password = util.PasswordHash(password, -1)
	db, err := util.GetDB()
	if err != nil {
		return util.FileNotFound
	}
	_, user_rev, err := db.Insert(&user)
	if err != nil {
		return util.MessagePage("Username already taken.", ctx)
	}
	//create a BlogData document for this user
	blogData := util.NewBlogData()
	blogData.ID = "BlogData_" + username
	_, blogData_rev, _ := db.Insert(&blogData)
	//if you can't add the user to the user list, delete the user
	if err != nil {
		db.Delete(user.ID, user_rev)
		db.Delete(blogData.ID, blogData_rev)
		return util.MessagePage("Error", ctx)
	}

	//if you can't add the user to the user list, delete the user
	if err != nil {
		db.Delete(user.ID, user_rev)
		return util.MessagePage("Error", ctx)
	}

	//return news of success
	if file, err := util.LoadTemplate("", "userCreated.html", ctx); err == nil {
		file = strings.Replace(file, "{{User.Name}}", username, -1)
		return file
	} else {
		return util.FileNotFound
	}
	return util.FileNotFound
}

func DeleteAccountPost(ctx *web.Context, val string) string {
	if db, err := util.GetDB(); err == nil {
		username := util.ReadUsername(ctx)
		var user User
		rev, err := db.Retrieve("User_"+username, &user)
		if err == nil && string(util.PasswordHash(ctx.Params["Password"], user.Password.Version).Hash) == string(user.Password.Hash) {
			if err := db.Delete("User_"+username, rev); err == nil {
				//delete the user's blog data
				bd := util.NewBlogData()
				rev, _ = db.Retrieve("BlogData_"+username, &bd)
				for i := 0; i < len(bd.Posts); i++ {
					post := posts.NewPost()
					postrev, _ := db.Retrieve(bd.Posts[i], &post)
					db.Delete(bd.Posts[i], postrev)
				}
				for i := 0; i < len(bd.Characters); i++ {
					char := char.NewCharacter()
					charrev, _ := db.Retrieve(bd.Characters[i], &char)
					db.Delete(bd.Characters[i], charrev)
				}
				db.Delete("BlogData_"+username, rev)
				//sign the user out
				if value, ok := util.ReadUserKey(ctx); ok {
					ctx.SetCookie("UserKey", value, -6000000)
					util.DeleteUserKey(value)
				}
				if username, ok := util.ReadCookie("Username", ctx); ok {
					ctx.SetCookie("Username", username, -6000000)
				}
				return util.MessagePage("The deed is done. You're dead to me now.", ctx)
			}
		}
	}
	return util.MessagePage("Operation failed, try again later.", ctx)
}
