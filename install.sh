if [ -n "$1" ]
then
	echo Installing as $1
	cd Code
	make link
	cd ../
	cp Code/main /bin/$1d

	rm -f /bin/$1
	rm -rf /var/www/$1
	cp -r ServeDir /var/www/$1
fi
