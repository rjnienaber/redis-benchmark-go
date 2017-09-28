#!/usr/bin/env bash

GOPATH=`pwd`
go build -o $GOPATH/bin/redis-benchmark main
