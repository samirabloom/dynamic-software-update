#
# Dynamic Software Upgrade Dockerfile
#

# Pull base image
FROM 127.0.0.1:5000/docker-go

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# Set up environment variables.
ENV DYNSOFTUP_HOME /home/goworld/src/github.com/samirabloom/software_upgrade/src

# Copy go files to container
WORKDIR /home/goworld/src/github.com/samirabloom/software_upgrade/src
ADD . /home/goworld/src/github.com/samirabloom/software_upgrade/src

# Expose ports
EXPOSE 8080

# Define default command
CMD ["go", "run", "docker_example.go"]