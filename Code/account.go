package main

import (
	"web"
)

func accountManagementGet(ctx *web.Context, val string) string {
	file, err := LoadTemplate("Accout Management", val, ctx)
	if err != nil {return fileNotFound}

	return file
}

func accountManagementPost(ctx *web.Context, val string) string {
	return ""
}

