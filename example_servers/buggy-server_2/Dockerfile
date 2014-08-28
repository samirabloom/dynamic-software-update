FROM ubuntu:12.04
MAINTAINER Petr Hosek "petrhosek@gmail.com"

RUN apt-get update && apt-get install -y \
	build-essential \
	autoconf \
	automake \
	libtool \
	subversion \
	libglib2.0-dev \
	libev-dev \
	libssl-dev \
	libpcre3-dev \
	zlib1g-dev \
	libbz2-dev \
	liblua5.1-dev \
	wget \
	libncurses5-dev

RUN wget http://www.lua.org/ftp/lua-5.1.4.tar.gz && \
    tar xvfz lua-5.1.4.tar.gz && \
    cd lua-5.1.4 && \
    make linux && \
    make install && \
    cp etc/lua.pc /usr/lib/pkgconfig/

RUN svn export -r 2148 svn://svn.lighttpd.net/lighttpd/branches/lighttpd-1.4.x/ /usr/src/lighttpd
WORKDIR /usr/src/lighttpd


VOLUME ["/var/www/html"]

WORKDIR /var/www/html

EXPOSE 80
CMD ["lighttpd", "-D", "-f", "/etc/lighttpd.conf"]
