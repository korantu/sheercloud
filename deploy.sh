#!/bin/bash
killall tool-launcher.sh
killall too
STAMP=`date +deploy%Y%m%d%H%M%S`
mkdir /tmp/$STAMP
cd /tmp/$STAMP
git init
git pull ~/git/sheercloud master
export GOPATH=/tmp/$STAMP/server
go build tool && cp ./tool ~/bin/tool 
tool-launcher.sh &
