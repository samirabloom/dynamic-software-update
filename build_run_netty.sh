#!/bin/bash

logLevel=$1
if [ -z "${logLevel}" ]; then logLevel="NOTICE"; fi

# note: using this instead of -a flag on go build to only rebuild local packages
echo "Cleaning previously built packages for current project"
rm -rf .pkg ./dynsoftup ./performance_log.csv
sleep 1

GOPATH=$PWD:$GOPATH
echo "Using GOROOT=${GOROOT}"
echo "Using GOPATH=${GOPATH}"

current_directory=$PWD

echo "Building example servers"
cd ./example_servers/netty_server/src/
javac -cp ../netty-all-4.0.21.Final.jar TestServer.java

echo "Running exmaple servers with 1034, 1035 and 1036"
java -cp .:../netty-all-4.0.21.Final.jar TestServer 1034 &
java -cp .:../netty-all-4.0.21.Final.jar TestServer 1035 &
java -cp .:../netty-all-4.0.21.Final.jar TestServer 1036 &

testServer1034Pid=`lsof -i:1034 | grep -v PID | awk '{print $2}'`
testServer1035Pid=`lsof -i:1035 | grep -v PID | awk '{print $2}'`
testServer1036Pid=`lsof -i:1036 | grep -v PID | awk '{print $2}'`

cd $current_directory

echo "Building project"
go build -o dynsoftup ./src/main_run.go

echo "Running proxy with logLevel ${logLevel}"
proxy -logLevel="${logLevel}" -configFile="config/config_script.json" &

trap "pkill dynsoftup; kill $testServer1034Pid $testServer1035Pid $testServer1036Pid" exit INT TERM

wait