#
# Dynamic Software Upgrade Go Server Dockerfile
#

# Pull base image
FROM samirabloom/docker-go

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# copy go files to container
ADD . .

# define default command
CMD ["go", "run", "go_server.go", "-port", "1025"]