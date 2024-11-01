#!/bin/sh
cd /srv/hugo
echo " === copy frontend ==="
cp -rf /app/static/build /srv/hugo/static/
cp -rf /app/data/manifest.json /srv/hugo/data/
echo " === generate pages ==="
if [ "$DO_NOT_MINIFY_HUGO" != "true" ]; then
  echo " === build with minify ==="
  hugo --minify
else
  echo " === build without minify ==="
  hugo
fi
/usr/local/bin/rss_generator
