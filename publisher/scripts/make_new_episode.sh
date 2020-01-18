#!/bin/sh

currdir=$(dirname $0)
cd ${currdir}/../
echo "current dir=`pwd`"

post=$(invoke print-next-episode-number 2>/dev/null)

echo "new post number=$post"
cd ../hugo

today=$(date +%Y-%m-%d)
hhmmss=$(date +%H:%M:%S)

outfile="./content/posts/podcast-$post.md"

echo '+++' >${outfile}
echo "title = \"Радио-Т $post\"" >>${outfile}
echo "date = \"${today}T${hhmmss}\"" >>${outfile}
echo 'categories = ["podcast"]' >>${outfile}
echo "image = \"https://radio-t.com/images/radio-t/rt$post.jpg\"" >>${outfile}
echo "filename = \"rt_podcast${post}\"" >>${outfile}
echo '+++' >>${outfile}
echo "" >>${outfile}
echo "![](https://radio-t.com/images/radio-t/rt${post}.jpg)" >>${outfile}
echo "" >>${outfile}

wget https://news.radio-t.com/api/v1/news/lastmd/12 -O /tmp/last-temi.tmp
cat /tmp/last-temi.tmp >>${outfile}
echo "- Темы наших слушателей" >>${outfile}
echo "" >>${outfile}

echo "*Спонсор этого выпуска [DigitalOcean](https://www.digitalocean.com)*
" >>${outfile}
echo "" >>${outfile}
echo "[аудио](https://cdn.radio-t.com/rt_podcast$post.mp3) • [лог чата](https://chat.radio-t.com/logs/radio-t-$post.html)" >>${outfile}
echo "<audio src=\"https://cdn.radio-t.com/rt_podcast$post.mp3\" preload=\"none\"></audio>" >>${outfile}

echo "new episode generated. File:"
echo "${outfile}"
