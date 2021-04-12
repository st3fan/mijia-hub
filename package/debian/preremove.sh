#!/bin/sh
set -e
if [ "$1" = remove ]; then
    /bin/systemctl stop mijia-hub
    /bin/systemctl disable mijia-hub
    /bin/systemctl daemon-reload
fi

