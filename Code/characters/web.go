/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package char

import (
	"web"
	"strconv"
	"strings"
	"libadv/util"
)

func ViewCharacterGet(ctx *web.Context, val string) string {
	charid, ok := ctx.Params["CharID"]
	if ok {
		db, err := util.GetDB()
		if err != nil {
			return util.FileNotFound
		}
		char := NewCharacter()
		db.Retrieve(charid, &char)
		if file, err := util.LoadTemplate(char.Name, "Character.html", ctx); err == nil {
			file = strings.Replace(file, "{{Name}}", char.Name, -1)
			file = strings.Replace(file, "{{Game}}", char.Game, -1)
			file = strings.Replace(file, "{{World}}", char.World, -1)
			file = strings.Replace(file, "{{Alligence}}", char.Alligiance, -1)
			file = strings.Replace(file, "{{Bio}}", char.Bio, -1)
			return file
		}
	}
	return util.FileNotFound
}

func AddCharacterPost(ctx *web.Context, val string) string {
	game, ok := ctx.Params["Game"]
	if ok {
		if !util.SignedIn(ctx) {
			return util.MessagePage("Please sign in.", ctx)
		}
		file, err := util.LoadTemplate("Add"+game+"Character", "CharacterEditor.html", ctx)

		if err == nil {
			file = strings.Replace(file, "{{Game}}", game, -1)
			file = strings.Replace(file, "{{CharacterID}}", "", -1)
			file = strings.Replace(file, "{{Name}}", "", -1)
			file = strings.Replace(file, "{{World}}", "", -1)
			file = strings.Replace(file, "{{Alligiance}}", "", -1)
			file = strings.Replace(file, "{{Bio}}", "", -1)
			return file
		}
	}
	return util.MessagePage("Operation failed, try again later.", ctx)
}

func EditCharacterGet(ctx *web.Context, val string) string {
	if db, err := util.GetDB(); err == nil {
		blog := util.NewBlogData()
		db.Retrieve("BlogData_"+util.ReadUsername(ctx), &blog)
		chars := "<option>----</option>\n"
		for i := 0; i < len(blog.Characters); i++ {
			char := NewCharacter()
			db.Retrieve(blog.Characters[i], &char)
			chars += "\t\t<option value=\"" + blog.Characters[i] + "\">" + char.Name + " (" + char.Game + " - " + char.World + ")</option>\n"
		}
		file, err := util.LoadTemplate("Edit Character", "EditCharacter.html", ctx)
		if err == nil {
			file = strings.Replace(file, "{{CharacterOptions}}", chars, -1)
			return file
		}
	}
	return util.FileNotFound
}

func EditCharacterPost(ctx *web.Context, val string) string {
	if db, err := util.GetDB(); err == nil {
		var char Character
		_, err = db.Retrieve(ctx.Params["CharacterID"], &char)
		if err == nil {
			if util.ReadUsername(ctx) == char.Owner {
				file, err := util.LoadTemplate("Editing "+char.Name, "CharacterEditor.html", ctx)
				if err == nil {
					file = strings.Replace(file, "{{CharacterID}}", ctx.Params["CharacterID"], -1)
					file = strings.Replace(file, "{{Name}}", char.Name, -1)
					file = strings.Replace(file, "{{Game}}", char.Game, -1)
					file = strings.Replace(file, "{{World}}", char.World, -1)
					file = strings.Replace(file, "{{Alligiance}}", char.Alligiance, -1)
					file = strings.Replace(file, "{{Bio}}", char.Bio, -1)
					return file
				}
			}
		}
	}
	return util.FileNotFound
}

func CharacterEditorPost(ctx *web.Context, val string) string {
	if util.SignedIn(ctx) {
		char := NewCharacter()
		char.Owner = util.ReadUsername(ctx)
		char.ID = ctx.Params["CharacterID"]
		char.Game = ctx.Params["Game"]
		char.Name = ctx.Params["Name"]
		char.World = ctx.Params["World"]
		char.Alligiance = ctx.Params["Alligiance"]
		char.Bio = ctx.Params["Bio"]
		if db, err := util.GetDB(); err == nil {
			blog := util.NewBlogData()
			db.Retrieve("BlogData_"+char.Owner, &blog)
			dummy := NewCharacter()
			rev, err := db.Retrieve(char.ID, &dummy)
			if err == nil {
				if dummy.Owner != char.Owner {
					return util.MessagePage("You are not authorized to edit this charater.", ctx)
				}
				char.Rev = rev
				db.Edit(&char)
				return util.MessagePage("Character updated.", ctx)
			} else {
				char.ID = "Character_" + strconv.Itoa(blog.CharacterIndex) + "_" + char.Owner
				db.Insert(&char)
				blog.CharacterIndex++
				blog.Characters = append(blog.Characters, char.ID)
				db.Edit(&blog)
				return util.MessagePage("Character created.", ctx)
			}
		}
	}
	return util.MessagePage("Operation failed, try again later.", ctx)
}
