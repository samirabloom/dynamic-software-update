#!/bin/bash

logLevel=$1
if [ -z "${logLevel}" ]; then logLevel="WARN"; fi

# note: using this instead of -a flag on go build to only rebuild local packages
echo "Cleaning previously built packages for current project"
rm -rf .pkg/* ./dynsoftup ./performance_log.csv
sleep 1

GOPATH=$PWD:$GOPATH
echo "Using GOROOT=${GOROOT}"
echo "Using GOPATH=${GOPATH}"

echo "Building project"
go build -o dynsoftup ./src/main_run.go

echo "Running main func with logLevel ${logLevel}"
./dynsoftup -logLevel="${logLevel}"