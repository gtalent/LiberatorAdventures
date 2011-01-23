include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
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
