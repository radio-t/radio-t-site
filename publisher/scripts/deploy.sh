#!/bin/sh
currdir=`dirname $0`
cd ${currdir}/../
echo "current dir=`pwd`"

echo "generates site"
num_before=`docker-compose run publisher print-next-episode-number 2>/dev/null`

cd ..
git pull
git add .
git commit -m "auto episode after $num_before" && git push
ssh master.radio-t.com "cd /srv/site.hugo && git pull && docker-compose run --rm hugo"

cd publisher
num_after=`docker-compose run --rm publisher print-next-episode-number 2>/dev/null`
ssh master.radio-t.com "docker exec -i gitter-bot /srv/gitter-rt-bot --super=Umputun --super=bobuk --super=ksenks --super=grayru --dbg --export-num=$num_before --export-path=/srv/html"

if [[ $num_before != $num_after ]]
then
  link=`docker-compose run --rm publisher print-last-rt-link`
  echo "will post new tweet for link $link"
  #./rt.tweet "радио-т $num_before $link #radiot"
fi

echo "remove articles"
http -a ${RT_NEWS_ADMIN} DELETE https://news.radio-t.com/api/v1/news/active/last/8
echo "Done"
