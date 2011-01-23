package main

import (
	"web"
)

func accountManagementGet(ctx *web.Context, val string) string {
	file, err := LoadTemplate("Accout Management", val, ctx)
	if err != nil {
		return fileNotFound
	}

	return file
}

func accountManagementPost(ctx *web.Context, val string) string {
	return ""
}

func deleteAccountPost(ctx *web.Context, val string) string {
	if db, err := getDB(); err == nil {
		username := readUsername(ctx)
		var user User
		rev, err := db.Retrieve("User_"+username, &user)
		if err == nil && ctx.Params["Password"] == user.Password {
			if err := db.Delete("User_"+username, rev); err == nil {
				return messagePage("The deed is done. You're dead to me now.", ctx)
			}
		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}
