#!/bin/bash

# Build Script for creating multiple architecture releases

# Requires:
# go get github.com/mitchellh/gox

## get version from self to include in file names
go build
VERSION=`wpsync --version | sed -e 's/wpsync v//'`

echo "Building $VERSION"
gox -osarch="linux/amd64 darwin/amd64 windows/amd64" -output "{{.Dir}}-$VERSION-{{.OS}}/{{.Dir}}"

for arch in linux darwin windows; do
    tar cf wpsync-$VERSION-$arch.tar wpsync-$VERSION-$arch
    gzip wpsync-$VERSION-$arch.tar
    rm -rf wpsync-$VERSION-$arch
done

