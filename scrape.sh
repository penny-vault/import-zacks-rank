#!/bin/sh
cd /usr/src/app
npm start
./import-zacks-rank Downloads/*.csv
