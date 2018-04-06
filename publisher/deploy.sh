#!/bin/sh
currdir=`dirname $0`
echo "current dir=$currdir"
cd ${currdir}

echo "generates site"
num_before=`utils/get-next-rt.py 2>/dev/null`

cd ..
git add .
git commit -m "auto episode after $num_before" && git push
ssh master.radio-t.com "cd /srv/site.hugo && git pull && docker-compose run --rm hugo"

cd publisher
num_after=`utils/get-next-rt.py 2>/dev/null`
ssh master.radio-t.com "docker exec -i gitter-bot /srv/gitter-rt-bot --super=Umputun --super=bobuk --super=ksenks --super=grayru --dbg --export-num=$num_before --export-path=/srv/html"

if [[ $num_before != $num_after ]]
then
  link=`utils/get-last-rt-link.py`
  echo "will post new tweet for link $link"
  #./rt.tweet "радио-т $num_before $link #radiot"
fi

echo "remove articles"
http -a ${RT_NEWS_ADMIN} DELETE https://news.radio-t.com/api/v1/news/active/last/8
echo "Done"
