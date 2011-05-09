/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package main

type MaterialSolution struct {
	Name	    string
	Requirement string
	Solution    string
	Quantity    uint64
}

type Schematic struct {
	ID             string "_id"
	Rev            string "_rev"
	Type           string
	Owner	  string
	Name      string
	Materials []MaterialSolution
}

func NewSchematic() Schematic {
	var schem Schematic
	schem.Materials = make([]MaterialSolution, 1)
	return schem
}

func (me *Schematic) AddMaterial(mat MaterialSolution) {
	me.Materials = append(me.Materials, mat)
}
