#!/bin/bash
killall tool
STAMP=`date +deploy%Y%m%d%H%M%S`
mkdir /tmp/$STAMP
cd /tmp/$STAMP
git init
git pull ~/git/sheercloud
go build tool
