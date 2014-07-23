#!/usr/bin/env bash

# == BUILD DOCKERFILE IN LOCAL DIRECTORY BASED ON GO BASE IMAGE IN LOCAL REGISTRY ==
echo
echo "============================================================================"
echo "build proxy image based on go base image (from public registry)"
echo "============================================================================"
echo

# build proxy image based on go base image (from public registry)
docker build -t samirabloom/proxy .

# make sure any containers that may exist from a previous run are stopped and removed
docker stop proxy 2> /dev/null
docker rm proxy 2> /dev/null

# run container mapping internal docker port 1234 to external ports 1234 (i.e. -p <external port>:<internal port>
#docker run -d --name proxy -p 1234:1234 -v /vagrant/config:/dynamic_software_update_config samirabloom/proxy
docker run --name proxy -p 1234:1234 -v /vagrant/config:/dynamic_software_update_config samirabloom/proxy

echo
echo "========================================"
echo "proxy docker is running successfully"
echo "========================================"
echo