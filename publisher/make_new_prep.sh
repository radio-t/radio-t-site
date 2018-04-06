#!/bin/sh

currdir=`dirname $0`
cd ${currdir}
echo "current dir=$currdir"

post=`utils/get-next-rt.py 2>/dev/null`

echo "new post number=$post"
cd ../hugo

today=$(date +%Y-%m-%d)
hhmmss=$(date +%H:%M:%S)

outfile="./content/posts/prep-$post.md"

echo '+++' > ${outfile}
echo "title = \"Темы для $post\"" >> ${outfile}
echo "date = \"${today}T${hhmmss}\"" >> ${outfile}
echo 'categories = ["prep"]' >> ${outfile}
echo '+++' >> ${outfile}

st3 ${outfile} &
