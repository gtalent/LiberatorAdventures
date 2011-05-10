include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	characters/data.go\
	characters/web.go\
	html/forms.go\
	html/html.go\
	html/image.go\
	html/lists.go\
	html/page.go\
	posts/data.go\
	posts/web.go\
	schematics/data.go\
	schematics/web.go\
	users/data.go\
	users/web.go\
	util/channel.go\
	util/misc.go\
	util/settings.go\
	web.go\
	data.go\
	main.go\

include $(GOROOT)/src/Make.pkg

link: package
	$(O)l -o main _go_.$(O)
run: link
	#./main -p -conf /usr/local/etc/LiberatorAdventuresd.conf
	./main -p
