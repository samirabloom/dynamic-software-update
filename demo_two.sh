#!/usr/bin/env bash

GOOGLE_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'

echo
echo "=============="
echo "OPENING CHROME"
echo "=============="
echo

read -n 2 -s

USER_DATA_DIR=google/user/data/`date +'%s'`
"$GOOGLE_CHROME" --user-data-dir=$USER_DATA_DIR --no-default-browser-check --no-first-run --disable-default-apps --window-position=0,0 "http://127.0.0.1:8080" &> /dev/null

echo
echo "================================================================="
echo "ADDING INSTANT UPGRADE FOR WORDPRESS & MYSQL IN DOCKER CONTAINERS"
echo "================================================================="
echo

read -n 2 -s

curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/demo_config_curl_docker_wordpress_INSTANT.json

echo
echo
echo "==============================================="
echo "UPGRADED WORDPRESS & MYSQL IN DOCKER CONTAINERS"
echo "==============================================="
echo

read -n 2 -s

echo
echo "=============="
echo "OPENING CHROME"
echo "=============="
echo

read -n 2 -s

USER_DATA_DIR=google/user/data/`date +'%s'`
"$GOOGLE_CHROME" --user-data-dir=$USER_DATA_DIR --no-default-browser-check --no-first-run --disable-default-apps --window-position=0,0 "http://127.0.0.1:8080" &> /dev/null

echo
echo "==================================="
echo "FAILING UPGRADE WORDPRESS CONTAINER"
echo "==================================="
echo

read -n 2 -s

echo "password: vagrant"
ssh vagrant@192.168.50.7 'sudo docker stop some-wordpress'

echo
echo "===================================="
echo "STOPPED UPGRADED WORDPRESS CONTAINER"
echo "===================================="
echo

read -n 2 -s

echo
echo "=============="
echo "OPENING CHROME"
echo "=============="
echo

read -n 2 -s

USER_DATA_DIR=google/user/data/`date +'%s'`
"$GOOGLE_CHROME" --user-data-dir=$USER_DATA_DIR --no-default-browser-check --no-first-run --disable-default-apps --window-position=0,0 "http://127.0.0.1:8080" &> /dev/null

echo
echo "===================================================================="
echo "ADDING CONCURRENT UPGRADE FOR WORDPRESS & MYSQL IN DOCKER CONTAINERS"
echo "===================================================================="
echo

read -n 2 -s

curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/demo_config_curl_docker_wordpress_CONCURRENT.json

echo
echo "=============="
echo "OPENING CHROME"
echo "=============="
echo

read -n 2 -s

USER_DATA_DIR=google/user/data/`date +'%s'`
"$GOOGLE_CHROME" --user-data-dir=$USER_DATA_DIR --no-default-browser-check --no-first-run --disable-default-apps --window-position=0,0 "http://127.0.0.1:8080" &> /dev/null

echo
echo "==================================="
echo "FAILING UPGRADE WORDPRESS CONTAINER"
echo "==================================="
echo

read -n 2 -s

echo "password: vagrant"
ssh vagrant@192.168.50.7 'sudo docker stop some-wordpress'

echo
echo "===================================="
echo "STOPPED UPGRADED WORDPRESS CONTAINER"
echo "===================================="
echo

read -n 2 -s

echo
echo "=============="
echo "OPENING CHROME"
echo "=============="
echo

read -n 2 -s

USER_DATA_DIR=google/user/data/`date +'%s'`
"$GOOGLE_CHROME" --user-data-dir=$USER_DATA_DIR --no-default-browser-check --no-first-run --disable-default-apps --window-position=0,0 "http://127.0.0.1:8080" &> /dev/null