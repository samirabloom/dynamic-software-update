#!/bin/bash
# Check that this script hasn't been run already
if [ ! -f /var/log/vmsetup ];
then
    # update apt-get and setup basic tools
    export DEBIAN_FRONTEND=noninteractive && \
    sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y build-essential software-properties-common libssl-dev curl git htop man unzip vim wget pkg-config mercurial bzr && \

	# install docker package
    apt-get -y install docker.io
    ln -sf /usr/bin/docker.io /usr/local/bin/docker
    sed -i '$acomplete -F _docker docker' /etc/bash_completion.d/docker.io
    chmod 777 /var/run/docker.sock

    # setup go environment variables
    echo "export PATH=/usr/local/go/bin:$PATH" >> /home/vagrant/.bash_profile
    echo "export GOROOT=/usr/local/go" >> /home/vagrant/.bash_profile
    echo "export GOPATH=/home/goworld" >> /home/vagrant/.bash_profile
    source /home/vagrant/.bash_profile

    # install go lang
    apt-get install -y build-essential software-properties-common curl git htop man unzip vim wget pkg-config mercurial bzr
    curl -s https://storage.googleapis.com/golang/go1.2.2.src.tar.gz | tar -v -C /usr/local -xz
    cd /usr/local/go/src && ./make.bash

    echo ===========
    echo   SUCCESS
    echo ===========
    echo
    echo now run: vagrant ssh docker
    echo

	# Place a marker to identify that this script has been run once already
	touch /var/log/vmsetup
fi


