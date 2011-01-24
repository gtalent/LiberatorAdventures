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

func deleteAccountPost(ctx *web.Context, val string) string {
	if db, err := getDB(); err == nil {
		username := readUsername(ctx)
		var user User
		rev, err := db.Retrieve("User_"+username, &user)
		if err == nil && ctx.Params["Password"] == user.Password {
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

