/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

//Used for the <ul> ta.
type UnorderedList struct {
	list *ComplexPageElement
}

//Creates a new UnorderedList.
func NewUnorderedList() *UnorderedList {
	l := new(UnorderedList)
	l.list = NewComplexPageElement()
	return l
}

func (me *UnorderedList) String() string {
	return "<ul>" + me.list.String() + "</ul"
}

//Adds the given string to the list as a list item.
func (me *UnorderedList) Add(item string) {
	me.list.Add(SimplePageElement("<li>" + item + "</li>"))
}

//Adds the given PageElement to the list as a list item.
func (me *UnorderedList) AddElement(item PageElement) {
	me.list.Add(SimplePageElement("<li>" + item.String() + "</li>"))
}
