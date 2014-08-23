#!/bin/bash
# Check that this script hasn't been run already
if [ ! -f /var/log/vmsetup ];
then
		
	# Install necessary software
	apt-get -y update
	
	# Install additional software
	sudo apt-get -y install vim libssl0.9.8
	
	wget http://packages.couchbase.com/releases/3.0.0-beta2/couchbase-server_3.0.0-beta2_x86_64_ubuntu_1004.deb -O couchbase-server_3.0.0-beta2_x86_64_ubuntu_1004.deb
	sudo dpkg -i couchbase-server_3.0.0-beta2_x86_64_ubuntu_1004.deb
	
	sleep 20
	
	# Remove document size limit
	sed -i 's/return getStringBytes(json) > self.docBytesLimit;/return false/g' /opt/couchbase/lib/ns_server/erlang/lib/ns_server/priv/public/js/documents.js

	echo 'Initializing cluster...'
	sudo /opt/couchbase/bin/couchbase-cli cluster-init -c 127.0.0.1:8091 --cluster-init-username=Administrator --cluster-init-password=password --cluster-init-ramsize=1024 -u Administrator -p password

	# install sample database
	curl 'http://127.0.0.1:8091/sampleBuckets/install' -H 'Authorization: Basic QWRtaW5pc3RyYXRvcjpwYXNzd29yZA==' --data '["beer-sample"]'
	
	# Place a marker to identify that this script has been run once already
	touch /var/log/vmsetup
fi
