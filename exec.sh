#!/bin/sh
echo " === generate pages ==="
cd /srv/hugo
hugo
/srv/hugo/generate_rss.py
