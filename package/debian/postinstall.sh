#!/bin/sh
set -e

if [ "$1" = configure ]; then
    /bin/systemctl daemon-reload
    /bin/systemctl enable mijia-hub
fi

