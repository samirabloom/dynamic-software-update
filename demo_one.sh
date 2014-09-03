#!/usr/bin/env bash

pkill proxy
echo "password: vagrant"
ssh vagrant@192.168.50.5 'sudo rm -rf /var/lib/mysql'
echo "password: vagrant"
ssh vagrant@192.168.50.7 'sudo rm -rf /var/lib/mysql'

echo
echo "================"
echo "BUILDING PROJECT"
echo "================"
echo

read -n 2 -s

make

echo
echo "=================="
echo "MAKE FILE COMPLETE"
echo "=================="
echo

read -n 2 -s

echo
echo "======================================================"
echo "RUNNING PROXY WITH WORDPRESS & MYSQL DOCKER CONTAINERS"
echo "======================================================"
echo

read -n 2 -s

proxy -logLevel="INFO" -configFile="config/demo_config_docker_wordpress.json"
