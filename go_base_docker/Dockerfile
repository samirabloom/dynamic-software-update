#
# Go Base Dockerfile
#

# Pull base image
FROM ubuntu:14.04

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# Set up environment variables.
ENV PATH /usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV GOPATH /home/goworld

# Install basic packages, go and compile go files
RUN \
  export DEBIAN_FRONTEND=noninteractive && \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get -y upgrade && \
  apt-get install -y build-essential && \
  apt-get install -y software-properties-common && \
  apt-get install -y curl git htop man unzip vim wget pkg-config && \ 
  curl -s https://storage.googleapis.com/golang/go1.2.2.src.tar.gz | tar -v -C /usr/local -xz && \
  cd /usr/local/go/src && ./make.bash

# Define default command to force RUN to execute
CMD ["bash"]