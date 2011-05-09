/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"strings"
	"os"
	"couch-go.googlecode.com/hg"
)

type DesignDoc struct {
	ID    string "_id"
	Rev   string "_rev"
	Lang  string "language"
	Views map[string]map[string]string "views"
}

func view(label, code string) map[string]map[string]string {
	view := make(map[string]map[string]string)
	view[label] = make(map[string]string)
	view[label]["map"] = code
	return view
}

//Gets the database connection.
func getDB() (couch.Database, os.Error) {
	return couch.NewDatabase(Settings.DatabaseAddress(), "5984", Settings.Database())
}

//Initializes the database by adding the design documents.
func initDB() bool {
	db, err := getDB()
	if err != nil {
		return false
	}
	_, _, err1 := db.Insert(design_users)
	_, _, err2 := db.Insert(design_posts)
	_, _, err3 := db.Insert(design_characters)
	return err1 == nil && err2 == nil && err3 == nil
}

type Document struct {
	ID             string "_id"
	Rev            string "_rev"
	Type           string
}

type BlogData struct {
	Document
	CharacterIndex int
	Characters     []string
	PostIndex      int
	Posts          []string
}

//Returns a new BlogData object by value.
func NewBlogData() BlogData {
	var data BlogData
	data.Type = "BlogData"
	return data
}

type Cookies struct {
	Document
	UserKeys map[string]string
}

func NewCookies() *Cookies {
	retval := new(Cookies)
	retval.UserKeys = make(map[string]string)
	return retval
}

//Unfinished
//Checks the username to make sure it does not contain any illegal characters.
//Returns a message if the name is bad, and a boolean indicating whether or not the name was ok
func checkUsername(username string) (string, bool) {
	//badChars := ";\\ \"\n\t={}()[]"
	if strings.Contains(username, ";") || strings.Contains(username, "\\") || strings.Contains(username, " ") || strings.Contains(username, "=") {
		return "Usernames may not contian the following characters: (A-Z, a-z, 0-9, _, -)", false
	}
	return "", true
}
