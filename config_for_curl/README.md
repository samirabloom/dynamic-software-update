# example docker configuration update
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_docker.json

# example error configuration update
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_error.json

# example instant upgrade
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_instant_upgrade.json

# example concurrent upgrade
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_concurrent_upgrade.json

# example lighttpd concurrent upgrade
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_ligttpd_bug.json

# example wordpress concurrent upgrade
curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_docker_wordpress.json

# example list configurations
curl -v -X GET 'http://127.0.0.1:9090/configuration/cluster/'

# example delete configuration
curl -v -X DELETE 'http://127.0.0.1:9090/configuration/cluster/48c2a4e6-31ba-11e4-bcdb-28cfe9158b63'

# demo steps

### terminal 1
-01 - ./build_run_go.sh INFO
-09 - view rollback in window

### browser
-02 - view wordpress and - setup with site title "WordPress 3.9.1" - go to homepage / without logging in
-06 - view wordpress and - setup with site title "WordPress 3.9.2" - go to homepage / without logging in
-08 - view wordpress and - see "WordPress 3.9.1" in title
-11 - view wordpress and SESSION timeout / upgrade

### terminal 2
-03 - curl 127.0.0.1:9090/configuration/cluster/
-04 - curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_docker_wordpress_INSTANT.json
-05 - curl 127.0.0.1:9090/configuration/cluster/
-07 - ssh vagrant@192.168.50.7 'sudo docker stop some-wordpress'   --- password vagrant
-10 - curl -v -X PUT 'http://127.0.0.1:9090/configuration/cluster' -H 'Content-Type: application/json' -d @config_for_curl/config_curl_docker_wordpress_SESSION.json
-12 - curl -v -X DELETE 'http://127.0.0.1:9090/configuration/cluster/fec3636f-32f8-11e4-aea1-28cfe9158b63'
