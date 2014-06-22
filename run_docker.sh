#!/usr/bin/env bash

boot2dockerIp=$(boot2docker ip 2>/dev/null)

# ensure boot2docker is up and running
boot2docker init
boot2docker up

if [ -z "${boot2dockerIp}" ]; 
then
    boot2dockerIp=127.0.0.1
else
	export DOCKER_HOST=tcp://$boot2dockerIp:2375	
fi
echo "======================================="
echo "Using boot2docker ip as: $boot2dockerIp"
echo "======================================="

# == CREATE REGISTRY ==
echo
echo "================"
echo "CREATE REGISTRY"
echo "================"
echo

# stop local registry
docker stop local_registry 2> /dev/null
docker rm local_registry 2> /dev/null

# create & start local docker registry
docker run -d --name local_registry -p 5000:5000 registry

# wait for docker registry to start up
sleep 60

# == ADD GO BASE IMAGE TO LOCAL REGISTRY ==
echo
echo "==================================="
echo "ADD GO BASE IMAGE TO LOCAL REGISTRY"
echo "==================================="
echo

# pull go base from public registry
docker run samirabloom/docker-go

# remove any existing tag
docker rmi $boot2dockerIp:5000/docker-go 2> /dev/null

# tag docker image
docker tag $(docker images -q samirabloom/docker-go) $boot2dockerIp:5000/docker-go

# push go base to local registry
docker push $boot2dockerIp:5000/docker-go

# == BUILD DOCKERFILE IN LOCAL DIRECTORY BASED ON GO BASE IMAGE IN LOCAL REGISTRY ==
echo
echo "============================================================================"
echo "BUILD DOCKERFILE IN LOCAL DIRECTORY BASED ON GO BASE IMAGE IN LOCAL REGISTRY"
echo "============================================================================"
echo

# build dynamic software update image based on go base image (in local registry)
docker build -t samirabloom/dynsoftup .

# make sure any containers that may exist from a previous run are stopped and removed
docker stop dynsoftup 2> /dev/null
docker rm dynsoftup 2> /dev/null

# run container mapping internal docker port 8080 to external ports 9090
docker run -d --name dynsoftup -p 9090:8080 samirabloom/dynsoftup

# remove any existing tag
docker rmi $boot2dockerIp:5000/dynsoftup 2> /dev/null

# tag docker image
docker tag $(docker images -q samirabloom/dynsoftup) $boot2dockerIp:5000/dynsoftup

# push dynsoftup to local registry
docker push $boot2dockerIp:5000/dynsoftup

# Note: boot2docker
# boot2docker ip - to get ip address
# boot2docker ssh - ssh to boot2docker box