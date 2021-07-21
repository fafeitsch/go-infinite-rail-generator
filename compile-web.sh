#!/bin/bash

cd ./web/app
ng build
mv dist/app ../../html
rm -r dist
cd ..