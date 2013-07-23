#!/bin/bash
echo Stopping the tool
killall tool-launcher.sh
killall tool
STAMP=`date +deploy%Y%m%d%H%M%S`
mkdir /tmp/$STAMP
cd /tmp/$STAMP
echo Cloning repo
git init
git pull ~/git/sheercloud master
echo Building the tool
export GOPATH=/tmp/$STAMP/server
go build tool && cp ./tool ~/bin/tool && echo Success!
echo Starting the tool
nohup tool-launcher.sh  > /dev/null 2> /dev/null < /dev/null &
echo Done
exit 0
