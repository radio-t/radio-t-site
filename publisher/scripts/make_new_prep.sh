#!/bin/sh

currdir=`dirname $0`
cd ${currdir}/../
echo "current dir=`pwd`"

post=`invoke print-next-episode-number 2>/dev/null`

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

echo "next episode prep generated. File:"
echo "../hugo/${outfile}"
