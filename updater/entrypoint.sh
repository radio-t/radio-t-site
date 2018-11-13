#!/bin/sh
echo "$(date) start updater"

while true
do
    ssh umputun@172.17.42.1 "/srv/site.hugo/updater/updater.sh"
    sleep 10
done
