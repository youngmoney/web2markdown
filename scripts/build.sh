#!/usr/bin/env bash

GOOS=darwin GOARCH=arm64 go build -o bin/web2markdown-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o bin/web2markdown-linux-amd64 .
