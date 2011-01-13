cd Code
make link
cd ../
cp Code/main /bin/LiberatorAdventuresd

rm -rf /var/www/LiberatorAdventures
cp -r ServeDir /var/www/LiberatorAdventures

