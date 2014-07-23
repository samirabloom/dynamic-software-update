#!/usr/bin/env bash

# build image
docker build -t samirabloom/example_server .

# make sure any previous containers removed

# port 1025
docker stop example_server_1025 2> /dev/null
docker rm example_server_1025 2> /dev/null
# port 1026
docker stop example_server_1026 2> /dev/null
docker rm example_server_1026 2> /dev/null
# port 1027
docker stop example_server_1027 2> /dev/null
docker rm example_server_1027 2> /dev/null

# run containers

# port 1025
docker run -d --name example_server_1025 -p 1025:1025 samirabloom/example_server go run example_server.go -port 1025
# port 1026
docker run -d --name example_server_1026 -p 1026:1026 samirabloom/example_server go run example_server.go -port 1026
# port 1027
docker run -d --name example_server_1027 -p 1027:1027 samirabloom/example_server go run example_server.go -port 1027