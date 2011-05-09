workspace=workspace
mkdir -p $workspace
git clone git://github.com/hoisie/web.go.git $workspace/web.go
cd $workspace/web.go
make install
cd ../../

hg clone https://couch-go.googlecode.com/hg/ $workspace/couch-go
cd $workspace/couch-go
make install
cd ../../

rm -rf $workspace
