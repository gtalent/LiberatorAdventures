/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

//Represents an HTML document.
type Page struct {
	Title, CSS string
	Elements *ComplexPageElement
}

//Returns a new Page object.
func NewPage(title string) Page {
	var page Page
	page.Title = title
	page.Elements = NewComplexPageElement()
	return page
}

//Adds the given element to this Page.
func (me *Page) Add(element PageElement) {
	me.Elements.Add(element)
}

//Adds the given string to this page.
func (me *Page) Put(text string) {
	me.Elements.Add(SimplePageElement(text))
}

//Adds the given string to this page.
func (me *Page) Putln(text string) {
	me.Elements.Add(SimplePageElement(text) + BR)
}

//Converts this Page into text to be served.
func (me *Page) String() (retval string) {
	if len(me.CSS) != 0 {
		retval = "<html><head><style type=\"text/css\">" + me.CSS + "</style><title>" + me.Title + "</title></head><body>"
	} else {
		retval = "<html><head><style type=\"text/css\">" + DefaultCSS + "</style><title>" + me.Title + "</title></head><body>"
	}
	retval += me.Elements.String()
	retval += "</body></html>"
	return
}
