#!/bin/bash
set -e

apk add go

go get github.com/cespare/reflex

go run github.com/cespare/reflex --decoration=fancy -r '\.go$' -s -- \
go run ./cmd/api --email $ML_EMAIL --course $ML_COURSE --password $ML_PASSWORD