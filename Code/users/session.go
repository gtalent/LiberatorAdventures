package main

import (
	"web"
	"rand"
	"strconv"
)

func signinPost(ctx *web.Context, val string) string {
	username := ctx.Params["Username"]
	password := ctx.Params["Password"]
	user := NewUser()
	if db, err := getDB(); err == nil {
		if _, err = db.Retrieve("User_"+username, &user); err == nil {
			if password == user.Password {
				num := rand.Int63()
				key := username + "_" + strconv.Itoa64(num)
				cookies.UserKeys[key] = username
				ctx.SetCookie("UserKey", key, 6000000)
				return messagePage("You are now signed in.", ctx)
			}
			return messagePage("Invalid username and password combination.", ctx)
		}
	}
	return messagePage("Could not access the database.", ctx)
}
