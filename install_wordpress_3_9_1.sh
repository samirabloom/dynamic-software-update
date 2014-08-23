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

    echo
    echo ===================================================================
    echo "  NEXT OPERATION TAKES >5 MINS TO DOWNLOAD MySQL DOCKER CONTAINER  "
    echo ===================================================================
    echo
    date

    docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=mysecretpassword -d mysql

    date

    echo
    echo =======================================================================
    echo "  NEXT OPERATION TAKES >5 MINS TO DOWNLOAD WordPress DOCKER CONTAINER  "
    echo =======================================================================
    echo
    date

    docker run --name some-wordpress -v /var/lib/mysql:/var/lib/mysql --link some-mysql:mysql -p 80:80 -d wordpress:3.9.1

    date

cat << EOF > ~/restart_wordpress.sh
#!/bin/bash

# stop all containers
docker ps -a | grep -v CONTAINER | awk '{print \$1}' | xargs docker stop

# delete all containers
docker ps -a | grep -v CONTAINER | awk '{print \$1}' | xargs docker rm

# run mysql container
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=mysecretpassword -d mysql

sleep 10

# run wordpress container
docker run --name some-wordpress --link some-mysql:mysql -p 80:80 -d wordpress:3.9.1
EOF
    chmod u+x ~/restart_wordpress.sh

	# Place a marker to identify that this script has been run once already
	touch /var/log/vmsetup
fi


