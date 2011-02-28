package main

type MaterialSolution struct {
	Name	    string
	Requirement string
	Solution    string
	Quantity    uint64
}

type Schematic struct {
	ID        string "_id"
	Rev       string "_rev"
	Object	  string
	Owner	  string
	Type      string
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
