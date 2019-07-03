#!/bin/sh
cd /srv/hugo
echo " === copy frontend ==="
cp -rf /app/static/build /srv/hugo/static/
cp -rf /app/data/manifest.json /srv/hugo/data/
echo " === generate pages ==="
hugo --minify
/srv/hugo/generate_rss.py
