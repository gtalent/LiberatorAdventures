include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	dbviews/characters.go\
	dbviews/posts.go\
	dbviews/users.go\
	dbviews/design.go\
	schematics/data.go\
	schematics/schematics.go\
	users/data.go\
	users/session.go\
	web.go\
	posts.go\
	account.go\
	characters.go\
	data.go\
	main.go\

include $(GOROOT)/src/Make.pkg

link: package
	6l -o main _go_.6
run: link
	./main
