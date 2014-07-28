#!/bin/bash

echo "stopping all containers"
docker ps -a | grep -v CONTAINER | awk '{print $1}' | xargs docker stop

echo "removing all containers"
docker ps -a | grep -v CONTAINER | awk '{print $1}' | xargs docker rm

echo "remove all images (except ubuntu or samirabloom/docker-go)"
docker images | grep -v IMAGE | grep -v ubuntu | grep -v samirabloom/docker-go | grep -v nginx | awk '{print $3}' | xargs docker rmi -f