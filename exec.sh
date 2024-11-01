#!/bin/sh
cd /srv/hugo
echo " === copy frontend ==="
cp -rf /app/static/build /srv/hugo/static/
cp -rf /app/data/manifest.json /srv/hugo/data/
echo " === generate pages ==="
if [ -z "$DO_NOT_MINIFY_HUGO" ]; then
  echo " === build with minify ==="
  hugo --minify
else
  echo " === build without minify ==="
  hugo
fi
/usr/local/bin/rss_generator
