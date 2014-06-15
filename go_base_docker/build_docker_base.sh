#!/usr/bin/env bash

# ensure boot2docker is up and running
boot2docker init
boot2docker up
export DOCKER_HOST=tcp://$(boot2docker ip 2>/dev/null):2375

# build base GO docker image from Dockerfile
docker build -t samirarabbanian/go-docker-base .