#!/bin/bash

rm -rf ./dist
/bin/bash ./compile-web.sh
mkdir dist
mv html ./dist

LANG=en_EN
version=$(git describe --tags --abbrev=0)
time=$(date)
echo $version $time
go build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=${version:1}'" -o ./dist/rail-generator cmd/main/rail-generator.go