#!/bin/sh

set -x

rm -f mijia-hub
env GOOS=linux GOARCH=arm GOARM=5 go build \
  && env NFPM_ARCH=armhf nfpm pkg --packager deb

rm -f mijia-hub
env GOOS=linux GOARCH=arm64 go build \
  && env NFPM_ARCH=arm64 nfpm pkg --packager deb

