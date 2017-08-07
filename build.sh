#!/usr/bin/env bash

docker run --rm -v "$PWD"/../../..:/go/src -w /go/src/github.com/peti2001/jira-time-log -e CGO_ENABLED=0 -e GOOS=linux golang:1.8 go build  -a -tags netgo -ldflags '-w'
docker build -t peti2001/jira-time-log .
docker push peti2001/jira-time-log

rm jira-time-log