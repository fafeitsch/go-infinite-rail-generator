#!/bin/bash

rm -rf ./dist
/bin/bash ./compile-web.sh
mkdir dist
mv html ./dist
go build -o ./dist/rail-generator cmd/rail-generator.go