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

echo "Building example servers"
go build -o example_server ./example_servers/go_server/go_server.go

echo "Running exmaple servers with 1034, 1035 and 1036"
./example_server -port="1034" &
./example_server -port="1035" &
./example_server -port="1036" &

echo "Building project"
go build -o dynsoftup ./src/main_run.go

echo "Running proxy with logLevel ${logLevel}"
./dynsoftup -logLevel="${logLevel}" -configFile="config/config_script.json" &

trap "pkill dynsoftup; pkill example_server" exit INT TERM

wait