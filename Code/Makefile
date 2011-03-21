include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	characters/data.go\
	characters/web.go\
	posts/data.go\
	posts/web.go\
	schematics/data.go\
	schematics/web.go\
	users/data.go\
	users/web.go\
	web.go\
	data.go\
	main.go\

include $(GOROOT)/src/Make.pkg

link: package
	$(O)l -o main _go_.6
run: link
	./main
