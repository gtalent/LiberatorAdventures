/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package schematics

import (
	"web"
	"strings"
	"strconv"
	"libadv/html"
	"libadv/util"
)

func ViewSchematicGet(ctx *web.Context, val string) string {
	if file, err := util.LoadTemplate("", "Schematic.html", ctx); err == nil {
		schemID := ctx.Params["SchemID"]
		db, err := util.GetDB()
		if err != nil {
			return util.FileNotFound
		}

		schem := NewSchematic()
		db.Retrieve(schemID, &schem)

		materials := html.NewUnorderedList()
		materialTemp, err := util.LoadFile("widgets/Material.html")
		if err != nil {
			return util.FileNotFound
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
	return util.FileNotFound
}
