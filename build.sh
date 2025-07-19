#!/bin/sh
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
echo "Building GopherTube with version: $LATEST_TAG"
go build -ldflags "-X main.version=$LATEST_TAG" -o gophertube main.go   