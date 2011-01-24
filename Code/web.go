package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"rand"
	"os"
	"web"
	"blinz/server"
)

const cookie string = "LiberatorAdventures"

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
	case "EditPost.html":
		return getEditPost(ctx, val)
	case "Account.html":
		return accountManagementGet(ctx, val)
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
	case "view/", "view":
		return viewPost(ctx, val)
	case "", "index.html", "index.htm":
		db, err := getDB()
		data, err := LoadTemplate("", "index.html", ctx)
		if err != nil {
			break
		}
		list := ""
		if users, err := db.Query("_design/users/_view/all", nil); err == nil {
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
	case "signin.html":
		if signedIn(ctx) {
			return messagePage("You're already signed in.", ctx)
		}
		retval, err := LoadTemplate("", val, ctx)
		if err != nil {
			break
		}
		return retval

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
	case "AddSWGEmuCharacter.html":
		return addSWGEmuCharacterPost(ctx, val)
	case "EditPost.html":
		return postEditPost(ctx, val)
	case "signin.html":
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
		break
	case "DeleteAccount.html":
		return deleteAccountPost(ctx, val)
		break
	case "CreateUser":
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
		user.Password = password
		db, err := getDB()
		if err != nil {
			break
		}
		_, user_rev, err := db.Insert(&user)
		if err != nil {
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
			break
		}
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
