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
		file, err := LoadTemplate("Add"+game+"Character", "Edit"+game+"Character.html", ctx)
		strings.Replace(file, "{{Character}}", "", -1)
		if err == nil {
			return file
		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}

func editCharacterGet(ctx *web.Context, val string) string {
	if db, err := getDB(); err == nil {
		blog := NewBlogData()
		db.Retrieve("BlogData_"+readUsername(ctx), &blog)
		chars := "<option>----</option>\n"
		for i := 0; i < len(blog.Characters); i++ {
			char := NewCharacter()
			db.Retrieve(blog.Characters[i], &char)
			chars += "\t\t<option value=\"" + blog.Characters[i] + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}
		file, err := LoadTemplate("Edit Character", "EditCharacter.html", ctx)
		if err == nil {
			file = strings.Replace(file, "{{CharacterOptions}}", chars, -1)
			return file
		}
	}
	return fileNotFound
}

func editCharacterPost(ctx *web.Context, val string) string {
	if db, err := getDB(); err == nil {
		var char Character
		_, err = db.Retrieve(ctx.Params["CharacterID"], &char)
		if err == nil {
			if readUsername(ctx) == char.Owner {
				file, err := LoadTemplate("Editing "+char.Name, "Edit"+char.Game+"Character.html", ctx)
				if err == nil {
					file = strings.Replace(file, "{{CharacterID}}", ctx.Params["CharacterID"], -1)
					file = strings.Replace(file, "{{Name}}", char.Name, -1)
					file = strings.Replace(file, "{{World}}", char.World, -1)
					file = strings.Replace(file, "{{Alligiance}}", char.Alligiance, -1)
					file = strings.Replace(file, "{{Bio}}", char.Bio, -1)
					return file
				}
			}
		}
	}
	return fileNotFound
}

func editSWGEmuCharacterPost(ctx *web.Context, val string) string {
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
			rev, err := db.Retrieve(char.ID, &dummy)
			if err == nil {
				char.Rev = rev
				db.Edit(&char)
				return messagePage("Character updated.", ctx)
			} else if dummy.Owner == readUsername(ctx) {
				char.ID = "Character_" + strconv.Itoa(blog.CharacterIndex) + "_" + char.Owner
				db.Insert(&char)
				blog.CharacterIndex++
				blog.Characters = append(blog.Characters, char.ID)
				db.Edit(&blog)
				return messagePage("Character created.", ctx)
			} else {
				return messagePage("You are not authorized to edit this charater.", ctx)
			}
		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}
