#!/bin/bash

rm -rf ./dist
/bin/bash ./compile-web.sh
mkdir dist
mv html ./dist

LANG=en_EN
version=$(git describe --tags --abbrev=0)
time=$(date)
go build -ldflags="-X 'github.com/fafeitsch/go-infinite-rail-generator/version.BuildTime=$time' -X 'github.com/fafeitsch/go-infinite-rail-generator/version.BuildVersion=${version:1}'" -o ./dist/rail-generator cmd/rail-generator.go