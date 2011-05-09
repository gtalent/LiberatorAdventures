/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"web"
	"rand"
	"strings"
	"strconv"
)

//session management

func signinPost(ctx *web.Context, val string) string {
	username := ctx.Params["Username"]
	password := ctx.Params["Password"]
	user := NewUser()
	if db, err := getDB(); err == nil {
		if _, err = db.Retrieve("User_"+username, &user); err == nil {
			if string(PasswordHash(password)) == string(user.Password) {
				num := rand.Int63()
				key := username + "_" + strconv.Itoa64(num)
				cookies.UserKeys[key] = username
				ctx.SetSecureCookie("UserKey", key, 6000000)
				return messagePage("You are now signed in.", ctx)
			}
			return messagePage("Invalid username and password combination.", ctx)
		}
		return messagePage("Username not found.", ctx)
	}
	return messagePage("Could not access the database.", ctx)
}

//account management

func accountManagementGet(ctx *web.Context, val string) string {
	file, err := LoadTemplate("Accout Management", val, ctx)
	if err != nil {
		return fileNotFound
	}
	return file
}

func createAccountPost(ctx *web.Context, val string) string {
	username := ctx.Params["Username"]
	email := ctx.Params["Email"]
	password := ctx.Params["Password"]
	password2 := ctx.Params["Password2"]
	if password != password2 {
		return messagePage("Passwords do not match.", ctx)
	}
	if len(password) < 6 {
		return messagePage("You're password must be at least 6 characters long.", ctx)
	}
	if strings.Contains(username, ";") || strings.Contains(username, "\\") || strings.Contains(username, " ") || strings.Contains(username, "=") {
		return messagePage("Usernames may not contian the following characters: \" \", \"=\", \"\\\", or \";\".", ctx)
	}
	user := NewUser()
	user.Username = username
	user.ID = "User_" + username
	user.Email = email
	user.Password = PasswordHash(password)
	db, err := getDB()
	if err != nil {
		return fileNotFound
	}
	_, user_rev, err := db.Insert(&user)
	if err != nil {
		out.Put(err.String())
		return messagePage("Username already taken.", ctx)
	}
	//create a BlogData document for this user
	blogData := NewBlogData()
	blogData.ID = "BlogData_" + username
	_, blogData_rev, _ := db.Insert(&blogData)
	//if you can't add the user to the user list, delete the user
	if err != nil {
		db.Delete(user.ID, user_rev)
		db.Delete(blogData.ID, blogData_rev)
		return messagePage("Error", ctx)
	}

	//if you can't add the user to the user list, delete the user
	if err != nil {
		db.Delete(user.ID, user_rev)
		return messagePage("Error", ctx)
	}

	//return news of success
	if file, err := LoadTemplate("", "userCreated.html", ctx); err == nil {
		file = strings.Replace(file, "{{User.Name}}", username, -1)
		return file
	} else {
		return fileNotFound
	}
	return fileNotFound
}

func deleteAccountPost(ctx *web.Context, val string) string {
	if db, err := getDB(); err == nil {
		username := readUsername(ctx)
		var user User
		rev, err := db.Retrieve("User_"+username, &user)
		if err == nil && string(PasswordHash(ctx.Params["Password"])) == string(user.Password) {
			if err := db.Delete("User_"+username, rev); err == nil {
				//delete the user's blog data
				bd := NewBlogData()
				rev, _ = db.Retrieve("BlogData_"+username, &bd)
				for i := 0; i < len(bd.Posts); i++ {
					post := NewPost()
					postrev, _ := db.Retrieve(bd.Posts[i], &post)
					db.Delete(bd.Posts[i], postrev)
				}
				for i := 0; i < len(bd.Characters); i++ {
					char := NewCharacter()
					charrev, _ := db.Retrieve(bd.Characters[i], &char)
					db.Delete(bd.Characters[i], charrev)
				}
				db.Delete("BlogData_"+username, rev)
				//sign the user out
				if value, ok := readUserKey(ctx); ok {
					ctx.SetCookie("UserKey", value, -6000000)
					cookies.UserKeys[value] = "", false
				}
				if username, ok := readCookie("Username", ctx); ok {
					ctx.SetCookie("Username", username, -6000000)
				}
				return messagePage("The deed is done. You're dead to me now.", ctx)
			}
		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}
