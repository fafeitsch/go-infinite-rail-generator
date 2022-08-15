#!/bin/bash

cd ./web/infinite-rail-generator-webapp
npm run build
mv dist ../../html
cd ..