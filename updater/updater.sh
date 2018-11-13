#!/bin/sh
# this script runs outside of container, on host

cd /srv/site.hugo

git fetch;
LOCAL=$(git rev-parse HEAD);
REMOTE=$(git rev-parse @{u});

if [ $LOCAL != $REMOTE ]; then
    echo "$(date) git update detected"
    git pull origin master
    docker-compose run --rm hugo
    echo "$(date) update completed"
fi