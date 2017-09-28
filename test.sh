#!/usr/bin/env bash

GOPATH=`pwd`
cd src/benchmark/tests
ginkgo watch
