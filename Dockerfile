#
# Dynamic Software Upgrade Dockerfile
#

# Pull base image
FROM samirarabbanian/go-docker-base

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# Set up environment variables.
ENV DYNSOFTUP_HOME /home/goworld/src/github.com/samirarabbanian/software_upgrade/src

# Copy go files to container
WORKDIR /home/goworld/src/github.com/samirarabbanian/software_upgrade/src
ADD . /home/goworld/src/github.com/samirarabbanian/software_upgrade/src

# Expose ports
EXPOSE 8080

# Define default command
CMD ["go", "run", "test.go"]