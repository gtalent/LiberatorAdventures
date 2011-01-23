package main

import (
	"web"
)

func addCharacterPost(ctx *web.Context, val string) string {
	game, ok := ctx.Params["Game"]
	if ok {
		if !signedIn(ctx) {
			return messagePage("Please sign in.", ctx)
		}
		file, err := LoadTemplate("Add"+game+"Character", "Add"+game+"Character.html", ctx)
		if err == nil {
			return file
		}
		if _, err := getDB(); err == nil {

		}
	}
	return messagePage("Operation failed, try again later.", ctx)
}
