#!/bin/sh
currdir=`dirname $0`
cd ${currdir}/../
echo "current dir=`pwd`"

echo "generates site"
num_before=`invoke print-next-episode-number 2>/dev/null`

cd ..
git pull
git add .
git commit -m "auto episode after $num_before" && git push
ssh umputun@master.radio-t.com "cd /srv/site.hugo && git pull && docker-compose run --rm hugo"

cd publisher
ssh umputun@master.radio-t.com "docker exec -i super-bot /srv/telegram-rt-bot --super=Umputun --super=bobuk --super=ksenks --super=grayru --dbg --export-num=$num_before --export-path=/srv/html"
num_after=`invoke print-next-episode-number 2>/dev/null`

if [[ $num_before != $num_after ]]
then
  link=`invoke publisher print-last-rt-link`
  echo "will post new tweet for link $link"
  #./rt.tweet "радио-т $num_before $link #radiot"
fi

echo "remove articles"
http -a ${RT_NEWS_ADMIN} DELETE https://news.radio-t.com/api/v1/news/active/last/8
echo "Done"
