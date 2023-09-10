#!/bin/sh

# build site on cloudflare pages
# build command: ./build-cf.sh
# output directory: public
# root path: /hugo

set -e
npm run build
hugo --minify
pip install pytoml mistune==0.8.4
./generate_rss.py --save-to=public
