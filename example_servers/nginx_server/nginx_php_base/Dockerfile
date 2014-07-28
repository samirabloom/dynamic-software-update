#
# nginx with PHP5 Dockerfile
#

# Pull base image
FROM dockerfile/ubuntu

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

RUN \
  export DEBIAN_FRONTEND=noninteractive && \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get update && \
  apt-get -y install php5-common php5-cli php5-fpm && \
  apt-get install -y nginx && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  echo "cgi.fix_pathinfo = 0;" >> /etc/php5/fpm/php.ini && \
  mkdir /var/www && \
  echo "<?php phpinfo(); ?>" > /var/www/index.php

ADD default /etc/nginx/sites-available/default

CMD php5-fpm -D && nginx
