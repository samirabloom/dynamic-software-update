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
	liblua5.1-dev

RUN svn export -r 2438 svn://svn.lighttpd.net/lighttpd/branches/lighttpd-1.4.x/ /usr/src/lighttpd
WORKDIR /usr/src/lighttpd

RUN ./autogen.sh
RUN ./configure \
	--prefix=/usr \
	--exec-prefix=/usr \
	--bindir=/usr/bin \
	--sbindir=/usr/sbin \
	--libexecdir=/usr/libexec \
	--sysconfdir=/etc \
	--sharedstatedir=/usr/com \
	--localstatedir=/var \
	--libdir=/usr/lib \
	--includedir=/usr/include \
	--datarootdir=/usr/share \
	--disable-ipv6 \
	--with-openssl \
	--with-pcre \
	--with-zlib \
	--with-bzip2 \
	--with-lua

RUN make -j"$(nproc)"
RUN make install

ADD lighttpd.conf /etc/lighttpd.conf

RUN mkdir -p /var/www/html
RUN echo "<html><body><div>hello world</div></body><html>" > /var/www/html/index.html
RUN chown -R www-data:www-data /var/www/html

WORKDIR /var/www/html

EXPOSE 80

RUN echo "lighttpd -D -f /etc/lighttpd.conf &" > ~/lighttpd.sh
RUN chmod u+x ~/lighttpd.sh

CMD ["lighttpd", "-D", "-f", "/etc/lighttpd.conf"]















