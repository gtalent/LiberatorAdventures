package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"rand"
	"os"
	"web"
	"couch-go.googlecode.com/hg"
	"blinz/server"
)

var cookies map[string]string = make(map[string]string)

const cookie string = "LiberatorAdventures"

var fileNotFound string = "File not found, perhaps it was taken by Tusken Raiders?"
var out *server.ChannelLine

//Gets the database connection.
func getDB() (couch.Database, os.Error) {
	return couch.NewDatabase(server.Settings.DatabaseAddress(), "5984", "liberator_adventures")
}

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
		_, ok = cookies[key]
		return ok
	}
	return false
}

func readUsername(ctx *web.Context) string {
	if key, ok := readUserKey(ctx); ok {
		if username, ok := cookies[key]; ok {return username}
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
	if file, err := LoadFile("message.html"); err == nil {
		file = strings.Replace(file, "{{TopBar}}", TopBar(ctx), -1)
		file = strings.Replace(file, "{{Message}}", message, -1)
		return file
	}
	return fileNotFound
}

func LoadFile(path string) (string, os.Error) {
	data, err := ioutil.ReadFile(server.Settings.WebRoot() + path)
	return string(data), err
}

func home(ctx *web.Context, val string) string {
	switch val {
	case "EditPost.html":
		return getEditPost(ctx, val)
	case "Logout":
		if value, ok := readUserKey(ctx); ok {
			ctx.SetCookie("UserKey", value, -6000000)
			cookies[value] = "", false
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
		data, err := LoadFile("index.html")
		if err != nil {
			break
		}
		users := new(UserList)
		list := "<ul>\n"
		if _, err = db.Retrieve("UserList", users); err == nil {
			size := len(users.Users)
			for i := 0; i < size; i++ {
				user := users.Users[i]
				list += "<il><a href=\"" + server.Settings.WebHome() + "view?user=" + user + "\">" + user + "</a></il><br>"
			}
		}
		list += "</ul>"
		data = strings.Replace(data, "{{TopBar}}", TopBar(ctx), -1)
		data = strings.Replace(data, "{{UserList}}", list, -1)
		return data
	case "signin.html":
		if signedIn(ctx) {
			return messagePage("You're already signed in.", ctx)
		}
		retval, err := LoadFile(val)
		if err != nil {
			break
		}
		if strings.HasSuffix(val, ".wgt") || strings.HasSuffix(val, ".html") {
			retval = strings.Replace(retval, "{{TopBar}}", TopBar(ctx), -1)
		}
		return retval

	default:
		retval, err := LoadFile(val)
		if err != nil {
			break
		}
		retval = strings.Replace(retval, "{{TopBar}}", TopBar(ctx), -1)
		return retval
	}
	return fileNotFound
}

func post(ctx *web.Context, val string) string {
	switch val {
	case "EditPost.html":
		return postEditPost(ctx, val)
	case "signin.html":
		username := ctx.Params["Username"]
		password := ctx.Params["Password"]
		user := new(User)
		if db, err := getDB(); err == nil {
			if _, err = db.Retrieve("User_"+username, user); err == nil {
				if password == user.Password {
					num := rand.Int63()
					key := username + "_" + strconv.Itoa64(num)
					cookies[key] = username
					ctx.SetCookie("UserKey", key, 6000000)
					return messagePage("You are now signed in.", ctx)
				}
				return messagePage("Invalid username and password combination.", ctx)
			}
		}
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
		user := new(User)
		user.Username = username
		user.ID = "User_" + username
		user.Email = email
		user.Password = password
		db, err := getDB()
		if err != nil {
			break
		}
		_, user_rev, err := db.Insert(user)
		if err != nil {
			return messagePage("Username already taken.", ctx)
		}
		//create a BlogData document for this user
		blogData := new(BlogData)
		blogData.ID = "BlogData_" + username
		db.Insert(blogData)
		users := new(UserList)
		_, err = db.Retrieve("UserList", users)
		//if you can't add the user to the user list, delete the user
		if err != nil {
			db.Delete(user.ID, user_rev)
			return messagePage("Error", ctx)
		}

		users.Users = append(users.Users, user.Username)
		_, err = db.Edit(users)

		//if you can't add the user to the user list, delete the user
		if err != nil {
			db.Delete(user.ID, user_rev)
			return messagePage("Error", ctx)
		}

		//return news of success
		if file, err := LoadFile("userCreated.html"); err == nil {
			file = strings.Replace(file, "{{TopBar}}", TopBar(ctx), -1)
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
	s.Get("/Liberator/(.*)", home)
	s.Post("/Liberator/(.*)", post)
	s.Run("0.0.0.0:" + strconv.Uitoa(server.Settings.WebPort()))
}
