#!/bin/sh

set -x

env GOOS=linux GOARCH=arm GOARM=5 go build
