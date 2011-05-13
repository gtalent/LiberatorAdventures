/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 *	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"strings"
	"strconv"
	"log"
	"os"
	"web"
	"libadv/char"
	"libadv/posts"
	"libadv/schematics"
	"libadv/users"
	"libadv/util"
)

func get(ctx *web.Context, val string) string {
	switch val {
	case "Account.html":
		return users.AccountManagementGet(ctx, val)
	case "Character.html":
		return char.ViewCharacterGet(ctx, val)
	case "EditCharacter.html":
		return char.EditCharacterGet(ctx, val)
	case "EditPost.html":
		return posts.GetEditPost(ctx, val)
	case "", "index.html", "index.htm":
		db, err := util.GetDB()
		if err != nil {
			return util.MessagePage("Cannot access database.", ctx)
		}
		data, err := util.LoadTemplate("", "index.html", ctx)
		if err != nil {
			break
		}
		list := ""
		if users, err := db.QueryIds("_design/users/_view/all", nil); err == nil {
			list = "<ul>\n"
			size := len(users)
			for i := 0; i < size; i++ {
				user := strings.SplitAfter(users[i], "User_", 2)[1]
				list += "\t<il><a href=\"" + "view?user=" + user + "\">" + user + "</a></il><br>\n"
			}
			list += "</ul>"
		}
		data = strings.Replace(data, "{{UserList}}", list, -1)
		return data
	case "signout.html":
		if value, ok := util.ReadUserKey(ctx); ok {
			ctx.SetSecureCookie("UserKey", value, -6000000)
			util.DeleteUserKey(value)
		}
		if username, ok := util.ReadCookie("Username", ctx); ok {
			ctx.SetSecureCookie("Username", username, -6000000)
		}
		return util.MessagePage("You're signed out.", ctx)
		break
	case "Schematic.html":
		return schematics.ViewSchematicGet(ctx, val)
	case "signin.html":
		if util.SignedIn(ctx) {
			return util.MessagePage("You're already signed in.", ctx)
		}
		retval, err := util.LoadTemplate("", val, ctx)
		if err != nil {
			break
		}
		return retval
	case "view/", "view":
		return posts.ViewPost(ctx, val)
	default:
		if strings.HasSuffix(val, ".html") {
			retval, err := util.LoadTemplate("", val, ctx)
			if err != nil {
				break
			}
			return retval
		}
		retval, err := util.LoadFile(val)
		if err != nil {
			break
		}
		if strings.HasSuffix(val, ".html") {
		} else if strings.HasSuffix(val, ".wgt") {
			topbar, _ := util.TopBar(ctx)
			retval = strings.Replace(retval, "{{TopBar}}", topbar, -1)
		}
		return retval
	}
	return util.FileNotFound
}

func post(ctx *web.Context, val string) string {
	switch val {
	case "AddCharacter.html":
		return char.AddCharacterPost(ctx, val)
	case "CharacterEditor.html":
		return char.CharacterEditorPost(ctx, val)
	case "CreateUser":
		return users.CreateAccountPost(ctx, val)
	case "DeleteAccount.html":
		return users.DeleteAccountPost(ctx, val)
	case "EditCharacter.html":
		return char.EditCharacterPost(ctx, val)
	case "EditPost.html":
		return posts.PostEditPost(ctx, val)
	case "signin.html":
		return users.SigninPost(ctx, val)
	}
	return util.FileNotFound
}

type dummy struct{}

func (me dummy) Write(p []byte) (n int, err os.Error) {
	util.WebOut.Put(string(p))
	return 0, nil
}

func RunWebServer(line *util.ChannelLine) {
	util.WebOut = line
	web.SetLogger(log.New(new(dummy), "", 0))
	web.Config.CookieSecret = util.Settings.CookieSecret()
	web.Get("/Liberator/(.*)", get)
	web.Post("/Liberator/(.*)", post)
	web.Run("0.0.0.0:" + strconv.Uitoa(util.Settings.WebPort()))
}
