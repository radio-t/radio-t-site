#!/bin/bash

currdir=$(dirname $0)
cd ${currdir}/../
echo "current dir=`pwd`"

export LANG="en_US.UTF-8"
fname=$(basename $1)

episode=$(echo $1 | sed -n 's/.*rt_podcast\(.*\)\.mp3/\1/p')
echo "!notif: Radio-T detected #${episode}"
invoke set-mp3-tags $1
echo "!notif: Radio-T tagged"

cd /episodes

echo "upload to radio-t.com"
echo "!notif: upload started"
scp $1 umputun@master.radio-t.com:/srv/master-node/var/media/${fname}

echo "remove old media files"
ssh master.radio-t.com "find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'"

echo "run ansible tasks"
ssh master.radio-t.com "docker exec -i ansible /srv/deploy_radiot.sh $episode"

echo "copy to hp-usrv archives"
echo "!notif: copy to hp-usrv (local) archives"
scp -P 2222 $1 umputun@archives.umputun.com:/data/archive.rucast.net/radio-t/media/

echo "upload to archive site"
scp $1 umputun@master.radio-t.com:/data/archive/radio-t/media/${fname}
ssh umputun@master.radio-t.com "chmod 644 /data/archive/radio-t/media/${fname}"

echo "all done for $fname"
echo "!notif: all done for $fname"
