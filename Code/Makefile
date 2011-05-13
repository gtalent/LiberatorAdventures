include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	web.go\
	main.go\

include $(GOROOT)/src/Make.pkg

install-pkgs: install-util install-html install-characters install-posts install-users install-schematics
install-util:
	gomake -C util install
install-html:
	gomake -C html install
install-characters: install-util
	gomake -C characters install
install-posts: install-util install-html
	gomake -C posts install
install-schematics: install-util
	gomake -C schematics install
install-users: install-util install-posts
	gomake -C users install

clean-pkgs:
	gomake -C util clean
	gomake -C html clean
	gomake -C characters clean
	gomake -C posts clean
	gomake -C schematics clean
	gomake -C users clean

full-clean: clean clean-pkgs
	rm -f main

link: install-pkgs package
	$(O)l -o main _go_.$(O)

run: link
	./main -p
