/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

import (
	"time"
	"container/vector"
)

const (
	add         = iota
	remove      = iota
	checkout    = iota
	checkin     = iota
	killRoutine = iota
)

type Record interface {
	Open(key string) bool
	Save() bool
}

type catalogItem struct {
	value   Record
	checkedOut bool
}

//Passed into the service goroutine to add and remove records.
type catalogOp struct {
	operation int8
	key       string
	record    Record
	//returns the value that was in the spot if there was one, the passed in value if a fail addition to the table
	retval chan interface{}
}

//Creates a new CatalogOp object.
func newCatalogOp() *catalogOp {
	op := new(catalogOp)
	op.retval = make(chan interface{})
	return op
}

//Used to find loaded records by their key
type Catalog struct {
	name    string
	records map[string]catalogItem
	edit    chan *catalogOp
	running bool
	//a list of keys for items to remove
	toRemove vector.StringVector
	template Record
}

//Creates and returns a new Catalog.
func NewCatalog(name string) *Catalog {
	catalog := new(Catalog)
	catalog.name = name
	catalog.records = make(map[string]catalogItem)
	catalog.edit = make(chan *catalogOp)
	catalog.running = true
	return catalog
}

//Call as a GoRoutine to serve out records when requested.
func (me *Catalog) Run(channel *ChannelLine) {
	go func() {
		op := <-me.edit
		for op.operation != killRoutine {
			switch op.operation {
			case add:
				val, ok := me.records[op.key]
				if !ok {
					val.checkedOut = false
					val.value = op.record
					me.records[op.key] = val
					op.retval <- true
				} else {
					op.retval <- false
				}
			case checkout:
				val, ok := me.records[op.key]
				if ok && !val.checkedOut {
					val.checkedOut = true
					me.records[op.key] = val
					op.retval <- val.value
				} else {
					op.retval <- nil
				}
			case checkin:
				val, ok := me.records[op.key]
				if ok {
					val.checkedOut = false
					op.retval <- val.value
					me.toRemove.Push(op.key)
				}
			case remove:
				retval, ok := me.records[op.key]
				if ok {
					me.records[op.key] = retval, false
					op.retval <- retval.value
				} else {
					op.retval <- nil
				}
			default:
				channel.Put(me.name + ": Undefined operation.")
			}
			op = <-me.edit
		}
		me.running = false
		channel.Put("Terminating " + me.name + " Catalog.")
	}()
	//the cleanup routine
	go func() {
		for me.running {
			time.Sleep(600000000000)
			op := newCatalogOp()
			for me.toRemove.Len() != 0 {
				op.key = me.toRemove.Pop()
				op.operation = remove
				me.edit <- op
			}
		}
	}()
}

//Ends this Catalog's Go routine.
func (me *Catalog) Stop() {
	op := newCatalogOp()
	op.operation = killRoutine
	me.edit <- op
}

//Adds the given value to the catalog if the spot for the given key is empty.
//Returns true if successful, false otherwise
func (me *Catalog) Add(key string, value Record) bool {
	op := newCatalogOp()
	op.record = value
	op.key = key
	op.operation = add
	me.edit <- op
	return (<-op.retval).(bool)
}

/*
  Checks in the value at the given key in the table if there is one.
  Takes:
 	key
*/
func (me *Catalog) Checkin(key string) {
	op := newCatalogOp()
	op.key = key
	op.operation = checkin
	me.edit <- op
}

/*
  Checks out the value at the given key in the table if there is one.
  Takes:
 	key
  Returns:
 	the value at the given key
 	true if it was there, false otherwise
*/
func (me *Catalog) Checkout(key string) (Record, bool) {
	op := newCatalogOp()
	op.key = key
	op.operation = checkout
	me.edit <- op
	retval := <-op.retval
	if retval == nil {
		return nil, false
	}
	return retval.(Record), true
}

/*
  Returns the value associated with the given key without checking it out.
  You don't have to, and should not, check it back in.
  Takes:
  	key
  Returns:
 	the value at the given key
 	true if it was there, false otherwise
*/
func (me *Catalog) Peek(key string) (Record, bool) {
	val, ok := me.records[key]
	return val.value.(Record), ok
}
