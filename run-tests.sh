#!/bin/sh

#sudo apt-get install libsnappy-dev or equivalent

# pull down the snappy tests, since they work
curl -s http://snappy-go.googlecode.com/hg/snappy/snappy_test.go |sed 's/^package snappy/package csnappy/' >snappy_test.go
go get
go build
go test -v
