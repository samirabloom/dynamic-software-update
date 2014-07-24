#
# Netty Test Server Dockerfile
#

# Pull base image
FROM dockerfile/java

# Maintainer details
MAINTAINER Samira Rabbanian "samira.rabanian@gmail.com"

# Set working directory
RUN mkdir -p /opt/netty
WORKDIR /opt/netty

# Download and compile Netty test server
RUN wget -O netty-all-4.0.15.Final.jar http://search.maven.org/remotecontent?filepath=io/netty/netty-all/4.0.15.Final/netty-all-4.0.15.Final.jar
RUN javac -cp netty-all-4.0.15.Final.jar TestServer.java

# Expose ports.
EXPOSE 8080

# Define default command.
CMD ["java", "-cp", ".:netty-all-4.0.15.Final.jar", "TestServer"]
