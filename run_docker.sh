#!/usr/bin/env bash

# ensure boot2docker is up and running
boot2docker init
boot2docker up
export DOCKER_HOST=tcp://$(boot2docker ip 2>/dev/null):2375

# dynamic software update docker image from Dockerfile
docker build -t samirarabbanian/dynsoftup .

# make sure any containers that may exist from a previous run are stopped and removed
docker stop dynsoftup 2> /dev/null
docker rm dynsoftup 2> /dev/null

# run container mapping internal docker port 8080 to external ports 9090
docker run -d --name dynsoftup -p 9090:8080 samirarabbanian/dynsoftup

# Note: boot2docker
# boot2docker ip - to get ip address
# boot2docker ssh - ssh to boot2docker box