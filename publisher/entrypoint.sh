#!/bin/sh
set -e

if [ $# -lt 1  ]; then
    exec ls -a
fi

exec "$@"
