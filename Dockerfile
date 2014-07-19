#
# Dynamic Software Upgrade Dockerfile
#

# Pull base image
FROM 127.0.0.1:5000/docker-go

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# Setup correct GOPATH
ENV GOPATH /home/goworld:/dynamic_software_update

# install dependencies
RUN go get code.google.com/p/go-uuid/uuid && \
    go get github.com/op/go-logging && \
    mkdir /dynamic_software_update

# copy go files to container
ADD src /dynamic_software_update/src
ADD src/docker_main.go /dynamic_software_update/

# setup working directory
WORKDIR /dynamic_software_update

# VOLUME /dynamic_software_update_config

# expose ports
EXPOSE 1234

# define default command
CMD ["go", "run", "main_run.go", "-logLevel", "INFO", "-configFile", "/dynamic_software_update_config/config.json"]