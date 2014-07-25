#!/usr/bin/env bash

# build image
docker build -t samirabloom/example_netty_server .

# make sure any previous containers removed

# port 8080
docker stop example_netty_server_8080 2> /dev/null
docker rm example_netty_server_8080 2> /dev/null

## port 8081
docker stop example_netty_server_8081 2> /dev/null
docker rm example_netty_server_8081 2> /dev/null

## port 8082
docker stop example_netty_server_8082 2> /dev/null
docker rm example_netty_server_8082 2> /dev/null

# run containers
## port 8080
docker run -d --name example_netty_server_8080 -p 8080:8080 samirabloom/example_netty_server java -cp .:netty-all-4.0.15.Final.jar TestServer 8080

## port 8081
docker run -d --name example_netty_server_8081 -p 8081:8081 samirabloom/example_netty_server java -cp .:netty-all-4.0.15.Final.jar TestServer 8081

## port 8082
docker run -d --name example_netty_server_8082 -p 8082:8082 samirabloom/example_netty_server java -cp .:netty-all-4.0.15.Final.jar TestServer 8082