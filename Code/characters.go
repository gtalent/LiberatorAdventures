package main

import (
	"web"
	"strconv"
	"strings"
)

func addCharacterPost(ctx *web.Context, val string) string {
	game, ok := ctx.Params["Game"]
	if ok {
		if !signedIn(ctx) {
			return messagePage("Please sign in.", ctx)
		}
		file, err := LoadTemplate("Add"+game+"Character", "Add"+game+"Character.html", ctx)
		strings.Replace(file, "{{Character}}", "", -1)
		if err == nil {
			return file
		}
		if _, err := getDB(); err == nil {

		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}

func addSWGEmuCharacterPost(ctx *web.Context, val string) string {
	if signedIn(ctx) {
		char := NewCharacter()
		char.Owner = readUsername(ctx)
		char.ID = ctx.Params["CharacterID"]
		char.Game = ctx.Params["Game"]
		char.Name = ctx.Params["Name"]
		char.World = ctx.Params["World"]
		char.Alligiance = ctx.Params["Alligiance"]
		char.Bio = ctx.Params["Bio"]
		if db, err := getDB(); err == nil {
			blog := NewBlogData()
			db.Retrieve("BlogData_"+char.Owner, &blog)
			dummy := NewCharacter()
			_, err = db.Retrieve("Character_" + strconv.Itoa(blog.CharacterIndex) + "_" + char.Owner, &dummy)
			if err == nil {
				db.Edit(&char)
				return messagePage("Character updated.", ctx)
			} else {
				char.ID = "Character_" + strconv.Itoa(blog.CharacterIndex) + "_" + char.Owner
				db.Insert(&char)
				blog.CharacterIndex++
				blog.Characters = append(blog.Characters, char.ID)
				db.Edit(&blog)
				return messagePage("Character created.", ctx)
			}
		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}
