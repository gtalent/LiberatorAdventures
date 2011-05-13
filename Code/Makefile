include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	web.go\
	main.go\

include $(GOROOT)/src/Make.pkg

lib-util:
	gomake -C util install
lib-html:
	gomake -C html install
mod-characters:
	gomake -C characters install
mod-posts:
	gomake -C posts install
mod-schematics:
	gomake -C schematics install
mod-users:
	gomake -C users install
install-mods: lib-util lib-html mod-characters mod-posts mod-schematics mod-users
link: install-mods package
	$(O)l -o main _go_.$(O)
run: link
	./main -p
