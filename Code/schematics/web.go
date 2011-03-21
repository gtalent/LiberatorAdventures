package main

import (
	"web"
	"strings"
	"strconv"
	"blinz/html"
)

func viewSchematicGet(ctx *web.Context, val string) string {
	if file, err := LoadTemplate("", "Schematic.html", ctx); err == nil {
		schemID := ctx.Params["SchemID"]
		db, err := getDB()
		if err != nil {
			return fileNotFound
		}

		schem := NewSchematic()
		db.Retrieve(schemID, &schem)

		materials := html.NewUnorderedList()
		materialTemp, err := LoadFile("widgets/Material.html")
		if err != nil {
			return fileNotFound
		}
		size := len(schem.Materials)
		for i := 0; i < size; i++ {
			material := strings.Replace(materialTemp, "{{Name}}", schem.Materials[i].Name, -1)
			material = strings.Replace(material, "{{Requirement}}", schem.Materials[i].Requirement, -1)
			material = strings.Replace(material, "{{Solution}}", schem.Materials[i].Solution, -1)
			material = strings.Replace(material, "{{Quantity}}", strconv.Uitoa64(schem.Materials[i].Quantity), -1)
		}

		file = strings.Replace(file, "{{Materials}}", materials.String(), -1)
		return file
	}
	return fileNotFound
}
