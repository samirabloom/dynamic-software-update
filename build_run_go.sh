#!/bin/bash

logLevel=$1
if [ -z "${logLevel}" ]; then logLevel="INFO"; fi

GOPATH=$PWD:$GOPATH
echo "Using GOROOT=${GOROOT}"
echo "Using GOPATH=${GOPATH}"

echo "Building example servers"
go build -o example_server ./example_servers/go_server/go_server.go

echo "Running exmaple servers with 1034, 1035 and 1036"
./example_server -port="1034" &
./example_server -port="1035" &
./example_server -port="1036" &
./example_server -port="1037" &
./example_server -port="1038" &
./example_server -port="1039" &
./example_server -port="1040" &
./example_server -port="1041" &
./example_server -port="1042" &

echo "Building project"
#go build -o dynsoftup ./src/main_run.go
make

echo "Running proxy with logLevel ${logLevel}"
#./dynsoftup -logLevel="${logLevel}" -configFile="config/config_script.json" &
proxy -logLevel="${logLevel}" -configFile="config/config_script.json" &

trap "pkill proxy; pkill example_server" exit INT TERM

