/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 *	http://www.opensource.org/licenses/bsd-license.php
 */
package util

import (
	"io/ioutil"
	"strings"
	"os"
	"web"
)

var FileNotFound string = "File not found, perhaps it was taken by Tusken Raiders?"
var WebOut *ChannelLine
var cookies *Cookies = NewCookies()

func SetUserKey(key, username string) {
	cookies.UserKeys[key] = username
}

func DeleteUserKey(key string) {
	cookies.UserKeys[key] = "", false
}

func GetUserKey(key string) string {
	return cookies.UserKeys[key]
}

//Reads the requested cookie from the given cookie list.
//Returns the desired cookie value if present, and an ok boolean value to indicate success or failure
func ReadCookie(cookie string, ctx *web.Context) (string, bool) {
	c, ok := ctx.GetSecureCookie(cookie)
	return c, ok
}

func ReadUserKey(ctx *web.Context) (string, bool) {
	return ReadCookie("UserKey", ctx)
}

func SignedIn(ctx *web.Context) bool {
	if key, ok := ReadUserKey(ctx); ok {
		_, ok = cookies.UserKeys[key]
		return ok
	}
	return false
}

func ReadUsername(ctx *web.Context) string {
	if key, ok := ReadUserKey(ctx); ok {
		if username, ok := cookies.UserKeys[key]; ok {
			return username
		}
	}
	return ""
}

//Loads the bar at the top of the page with the title and session management links.
func TopBar(ctx *web.Context) (string, os.Error) {
	_, signedin := ReadUserKey(ctx)
	retval, err := LoadFile("TopBar.wgt")
	if err != nil {
		return "TopBar not found.", err
	}
	sessionManager := ""
	postManager := ""
	if !signedin {
		if file, err := LoadFile("SessionManagerAnon.wgt"); err == nil {
			sessionManager = file
		}
	} else {
		if file, err := LoadFile("SessionManager.wgt"); err == nil {
			sessionManager = file
		}
		if file, err := LoadFile("PostManagement.wgt"); err == nil {
			postManager = file
		}
	}
	retval = strings.Replace(retval, "{{WebHome}}", Settings.WebHome(), -1)
	retval = strings.Replace(retval, "{{SessionManager}}", sessionManager, -1)
	retval = strings.Replace(retval, "{{PostManagement}}", postManager, -1)
	return retval, nil
}

func PostDiv() string {
	bytes, err := ioutil.ReadFile(Settings.WebRoot() + "post.wgt")
	if err != nil {
		return "<div><h3>{{Title}}</h3><br>{{Content}}</div>"
	}
	return string(bytes)
}

func MessagePage(message string, ctx *web.Context) string {
	if file, err := LoadTemplate("", "message.html", ctx); err == nil {
		file = strings.Replace(file, "{{Message}}", message, -1)
		return file
	}
	return FileNotFound
}

func LoadFile(path string) (string, os.Error) {
	data, err := ioutil.ReadFile(Settings.WebRoot() + path)
	return string(data), err
}

//Load the template file and fills in the body with the contents of the file at the given path.
func LoadTemplate(subTitle, bodyPath string, ctx *web.Context) (string, os.Error) {
	data, err := LoadFile("template.html")
	if err != nil {
		return FileNotFound, err
	}
	if len(subTitle) != 0 {
		subTitle = " - " + subTitle
	}
	body, err := LoadFile(bodyPath)
	data = strings.Replace(data, "{{SubTitle}}", subTitle, -1)
	topbar, err := TopBar(ctx)
	data = strings.Replace(data, "{{TopBar}}", topbar, -1)
	data = strings.Replace(data, "{{Body}}", body, -1)
	return data, err
}
