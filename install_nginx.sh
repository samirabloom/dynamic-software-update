#!/bin/bash
# Check that this script hasn't been run already
if [ ! -f /var/log/vmsetup ];
then
    export DEBIAN_FRONTEND=noninteractive
    sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list
    apt-get update
    # apt-get -y upgrade
    apt-get install -y apache2 apache2-utils apache2-dev libpcre3 libpcre3-dev zlib1g-dev libxml2-dev liblua5.2-dev libcurl4-openssl-dev libyajl-dev
    cp /vagrant/{nginx-1.6.1.tar.gz,modsecurity-2.8.0.tar.gz,openssl-1.0.1h.tar.gz} /usr/local/src
    cd /usr/local/src
    tar -zxvf nginx-1.6.1.tar.gz
    tar -zxvf modsecurity-2.8.0.tar.gz
    tar -zxvf openssl-1.0.1h.tar.gz
    cd modsecurity-2.8.0
    ./configure
    make
    cd ../nginx-1.6.1
    ./configure --add-module=/usr/local/src/modsecurity-2.8.0/nginx/modsecurity --with-openssl=/usr/local/src/openssl-1.0.1h --with-http_ssl_module --with-debug
    make
fi


