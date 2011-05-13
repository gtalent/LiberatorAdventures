/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package html

import (
	"strconv"
)

var DefaultCSS string
var HR SimplePageElement = SimplePageElement("<hr>")
var BR SimplePageElement = SimplePageElement("<br>")

//An interface to represent an HTML page element.
type PageElement interface {
	String() string
}

//A PageElement for holding other page elements.
type SimplePageElement string


func (me SimplePageElement) String() string {
	return string(me)
}

//A PageElement for holding other page elements.
type ComplexPageElement struct {
	elements []PageElement
}

func NewComplexPageElement() (retval *ComplexPageElement) {
	retval = new(ComplexPageElement)
	retval.elements = make([]PageElement, 0)
	return retval
}

//Returns a string representing all the PageElements in this string.
func (me *ComplexPageElement) String() string {
	retval := ""
	for i := 0; i < len(me.elements); i++ {
		retval += me.elements[i].String()
	}
	return retval
}

//The number of elements in this ComplexPageElement.
func (me *ComplexPageElement) Size() int {
	return len(me.elements)
}

//Adds the given PageElement to this ComplexPageElement.
func (me *ComplexPageElement) Add(element PageElement) {
	me.elements = append(me.elements, element)
}

//Used to represent Cells in HTML tables.
type Cell string

//Returns this Cell formatted properly for insertion in an HTML table.
func (me Cell) String() string {
	return "<td>" + string(me) + "</td>"
}

//Used to represent a row in an HTML table.
type Row struct {
	Cells *ComplexPageElement
}

func NewRow() (row Row) {
	row.Cells = NewComplexPageElement()
	return
}

//Adds the given cell for display in this row.
func (me Row) AddCell(cell Cell) {
	me.Cells.Add(cell)
}

//Returns this row as displayable HTML.
func (me Row) String() string {
	return "<tr>" + me.Cells.String() + "</tr>"
}

//Used to represent HTML tables.
type Table struct {
	Alignment string
	Rows *ComplexPageElement
}

//Returns a new Table object.
func NewTable() Table {
	var table Table
	table.Rows = NewComplexPageElement()
	table.Alignment = "left"
	return table
}

//Adds the given Row to this Table.
func (me Table) AddRow(row Row) {
	me.Rows.Add(row)
}

//Returns the text of this Table ready to be inserted in an HTML document.
func (me Table) String() string {
	return "<table>" + me.Rows.String() + "</table>"
}

//Places the given text in H1 header brackets.
func H1(text string) SimplePageElement {
	return SimplePageElement("<h1>" + text + "</h1>")
}

//Places the given text in H2 header brackets.
func H2(text string) SimplePageElement {
	return SimplePageElement("<h2>" + text + "</h2>")
}

//Places the given text in H3 header brackets.
func H3(text string) SimplePageElement {
	return SimplePageElement("<h3>" + text + "</h3>")
}

//Writes out an attribute field for a tag.
func Attribute(field, value string) string {
	return " " + field + "=\"" + value + "\" "
}

//Writes out an attribute field for a tag.
func AttributeInt(field string, value int) string {
	return " " + field + "=\"" + strconv.Itoa(value) + "\" "
}

func TextLink(text, address string) SimplePageElement {
	return SimplePageElement("<a href=\"" + address + "\">" + text + "</a>")
}
