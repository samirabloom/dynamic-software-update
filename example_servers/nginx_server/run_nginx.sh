#!/usr/bin/env bash

# make sure any previous containers removed

## port 8085
docker stop example_nginx_server_8085 2> /dev/null
docker rm example_nginx_server_8085 2> /dev/null

## port 8086
docker stop example_nginx_server_8086 2> /dev/null
docker rm example_nginx_server_8086 2> /dev/null

## port 8087
docker stop example_nginx_server_8087 2> /dev/null
docker rm example_nginx_server_8087 2> /dev/null

# run containers

## port 8085
docker run --name example_nginx_server_8085 -p 8085:80 -v $PWD:/usr/local/nginx/html:ro -d nginx

## port 8086
docker run --name example_nginx_server_8086 -p 8086:80 -v $PWD:/usr/local/nginx/html:ro -d nginx

## port 8087
docker run --name example_nginx_server_8087 -p 8087:80 -v $PWD:/usr/local/nginx/html:ro -d nginx
