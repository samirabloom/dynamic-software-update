#!/usr/bin/env bash

# build image
docker build -t samirabloom/example_netty_server .

# make sure any previous containers removed

# port 1025
docker stop example_netty_server_8080 2> /dev/null
docker rm example_netty_server_8080 2> /dev/null
## port 1026
#docker stop example_go_server_8081 2> /dev/null
#docker rm example_go_server_8081 2> /dev/null
## port 1027
#docker stop example_go_server_8082 2> /dev/null
#docker rm example_go_server_8082 2> /dev/null

# run containers

# port 8080
docker run -d --name example_netty_server_8080 -p 8080:8080 samirabloom/example_netty_server
## port 1026
#docker run -d --name example_go_server_1026 -p 8081:8081 samirabloom/example_go_server go run go_server.go -port 8081
## port 1027
#docker run -d --name example_go_server_1027 -p 8082:8082 samirabloom/example_go_server go run go_server.go -port 8082