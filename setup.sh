workspace=workspace
mkdir -p $workspace
git clone git://github.com/hoisie/web.go.git $workspace/web.go
cd $workspace/web.go
gomake install
cd ../../

hg clone https://gtalent2@gtalent2-libadv.googlecode.com/hg/ $workspace/couch-go
cd $workspace/couch-go
gomake install
cd ../../

rm -rf $workspace
