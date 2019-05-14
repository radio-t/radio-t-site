#!/bin/sh
# this script runs outside of container, on host

cd /srv/site.hugo

git fetch;
LOCAL=$(git rev-parse HEAD);
REMOTE=$(git rev-parse @{u});

if [ $LOCAL != $REMOTE ]; then
    sleep 5
    echo "$(date) git update detected"
    git pull origin master
    docker-compose build hugo
    docker-compose run --rm hugo
    echo "$(date) update completed"
fi
