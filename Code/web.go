package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"os"
	"web"
	"blinz/server"
)


var fileNotFound string = "File not found, perhaps it was taken by Tusken Raiders?"
var out *server.ChannelLine

//Returns the given cookie list as map
func readCookies(ctx *web.Context) map[string]string {
	cookies := ctx.Headers["Cookie"]
	list := strings.Split(cookies, "; ", -1)
	size := len(list)
	m := make(map[string]string)
	for i := 0; i < size; i++ {
		pair := strings.Split(list[i], "=", -1)
		m[pair[0]] = pair[1]
	}
	return m
}

//Reads the requested cookie from the given cookie list.
//Returns the desired cookie value if present, and an ok boolean value to indicate success or failure
func readCookie(cookie string, ctx *web.Context) (string, bool) {
	cookies := ctx.Headers["Cookie"]
	list := strings.Split(cookies, "; ", -1)
	size := len(list)
	for i := 0; i < size; i++ {
		pair := strings.Split(list[i], "=", -1)
		if pair[0] == cookie {
			return pair[1], true
		}
	}
	return "", false
}

func readUserKey(ctx *web.Context) (string, bool) {
	return readCookie("UserKey", ctx)
}

func signedIn(ctx *web.Context) bool {
	if key, ok := readUserKey(ctx); ok {
		_, ok = cookies.UserKeys[key]
		return ok
	}
	return false
}

func readUsername(ctx *web.Context) string {
	if key, ok := readUserKey(ctx); ok {
		if username, ok := cookies.UserKeys[key]; ok {
			return username
		}
	}
	return ""
}

func TopBar(ctx *web.Context) string {
	_, signedin := readUserKey(ctx)
	retval, err := LoadFile("TopBar.wgt")
	if err != nil {
		return "TopBar not found."
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
	retval = strings.Replace(retval, "{{WebHome}}", server.Settings.WebHome(), -1)
	retval = strings.Replace(retval, "{{SessionManager}}", sessionManager, -1)
	retval = strings.Replace(retval, "{{PostManagement}}", postManager, -1)
	return retval
}

func postDiv() string {
	bytes, err := ioutil.ReadFile(server.Settings.WebRoot() + "post.wgt")
	if err != nil {
		return "<div><h3>{{Title}}</h3><br>{{Content}}</div>"
	}
	return string(bytes)
}

func messagePage(message string, ctx *web.Context) string {
	if file, err := LoadTemplate("", "message.html", ctx); err == nil {
		file = strings.Replace(file, "{{Message}}", message, -1)
		return file
	}
	return fileNotFound
}

func LoadFile(path string) (string, os.Error) {
	data, err := ioutil.ReadFile(server.Settings.WebRoot() + path)
	return string(data), err
}

//Load the template file and fills in the body with the contents of the file at the given path.
func LoadTemplate(subTitle, bodyPath string, ctx *web.Context) (string, os.Error) {
	data, err := LoadFile("template.html")
	if err != nil {
		return fileNotFound, err
	}
	if len(subTitle) != 0 {
		subTitle = " - " + subTitle
	}
	body, err := LoadFile(bodyPath)
	data = strings.Replace(data, "{{SubTitle}}", subTitle, -1)
	data = strings.Replace(data, "{{TopBar}}", TopBar(ctx), -1)
	data = strings.Replace(data, "{{Body}}", body, -1)
	return data, err
}

func get(ctx *web.Context, val string) string {
	switch val {
	case "Account.html":
		return accountManagementGet(ctx, val)
	case "Character.html":
		return viewCharacterGet(ctx, val)
	case "EditCharacter.html":
		return editCharacterGet(ctx, val)
	case "EditPost.html":
		return getEditPost(ctx, val)
	case "", "index.html", "index.htm":
		db, err := getDB()
		data, err := LoadTemplate("", "index.html", ctx)
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
	case "Logout":
		if value, ok := readUserKey(ctx); ok {
			ctx.SetCookie("UserKey", value, -6000000)
			cookies.UserKeys[value] = "", false
		}
		if username, ok := readCookie("Username", ctx); ok {
			ctx.SetCookie("Username", username, -6000000)
		}
		return messagePage("You're signed out.", ctx)
		break
	case "Schematic.html":
		return viewSchematicGet(ctx, val)
	case "signin.html":
		if signedIn(ctx) {
			return messagePage("You're already signed in.", ctx)
		}
		retval, err := LoadTemplate("", val, ctx)
		if err != nil {
			break
		}
		return retval
	case "view/", "view":
		return viewPost(ctx, val)
	default:
		if strings.HasSuffix(val, ".html") {
			retval, err := LoadTemplate("", val, ctx)
			if err != nil {
				break
			}
			return retval
		}
		retval, err := LoadFile(val)
		if err != nil {
			break
		}
		if strings.HasSuffix(val, ".html") {
		} else if strings.HasSuffix(val, ".wgt") {
			retval = strings.Replace(retval, "{{TopBar}}", TopBar(ctx), -1)
		}
		return retval
	}
	return fileNotFound
}

func post(ctx *web.Context, val string) string {
	switch val {
	case "AddCharacter.html":
		return addCharacterPost(ctx, val)
	case "CharacterEditor.html":
		return characterEditorPost(ctx, val)
	case "DeleteAccount.html":
		return deleteAccountPost(ctx, val)
	case "EditCharacter.html":
		return editCharacterPost(ctx, val)
	case "EditPost.html":
		return postEditPost(ctx, val)
	case "signin.html":
		return signinPost(ctx, val)
	case "CreateUser":
		return createAccountPost(ctx, val)
	}
	return fileNotFound
}

type dummy struct{}

func (me dummy) Write(p []byte) (n int, err os.Error) {
	return 0, nil
}

func RunWebServer(line *server.ChannelLine) {
	out = line
	var s web.Server
	s.Logger = log.New(new(dummy), "", 0)
	s.Get("/Liberator/(.*)", get)
	s.Post("/Liberator/(.*)", post)
	s.Run("0.0.0.0:" + strconv.Uitoa(server.Settings.WebPort()))
}
