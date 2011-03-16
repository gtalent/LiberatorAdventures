package main

type designDoc struct {
	ID    string "_id"
	Rev   string "_rev"
	Lang  string "language"
	Views map[string]map[string]string "views"
}

func view(label, code string) map[string]map[string]string {
	view := make(map[string]map[string]string)
	view[label] = make(map[string]string)
	view[label]["map"] = code
	return view
}
