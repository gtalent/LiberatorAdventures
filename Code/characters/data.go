package main 

var design_characters designDoc = designDoc{ID:   "_design/characters", Lang: "javascript",
	Views: view("all", "function(doc) { if (doc.Type == 'Character')  emit(doc.Name, doc) }")}

type Character struct {
	ID                                        string "_id"
	Rev                                       string "_rev"
	Type                                      string
	Game, Name, World, Alligiance, Bio, Owner string
}

func NewCharacter() Character {
	var data Character
	data.Type = "Character"
	return data
}

