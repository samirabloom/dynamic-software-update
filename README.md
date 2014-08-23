dynamic-software-update
=======================

# Installation

1. git clone https://github.com/samirabloom/dynamic-software-update
2. make

This will install the proxy to the `PATH` by adding it to the `/usr/local/bin` directory

# Usage

The proxy runs from the command line with the following options:

```bash
Usage of proxy:
  -configFile="./config.json": Set the location of the configuration file that should contain configuration to start the proxy,
                               for example:
                                           {
                                               "proxy": {
                                                   "port": 1235
                                               },
                                               "configService": {
                                                   "port": 9090
                                               },
                                               "cluster": {
                                                   "servers":[
                                                       {"ip": "127.0.0.1", "port": 1034},
                                                       {"ip": "127.0.0.1", "port": 1035}
                                                   ],
                                                   "version": 1.0
                                               }
                                           }

  -logLevel="WARN": Set the log level as "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO" or "DEBUG"
  
  -h: Displays this message
```

For example:

```bash
proxy -logLevel=INFO -configFile="config/config_script.json"
```

 
## wrk testing

### 1. Build wrk
```bash
cd ~/git/
git clone https://github.com/wg/wrk.git
cd wrk
make
```
 
### 2. Test Server Directly
 
```bash
./wrk -t200 -c200 -d10 --latency http://127.0.0.1:1024
```

### 3. Test Server Via Proxy

```bash
./wrk -t200 -c200 -d10 --latency http://127.0.0.1:1235
```

## ApacheBench
 
### 1. Test Server Directly
 
```bash
ab -n 10000 -c 100 http://127.0.0.1:1024/
```

### 2. Test Server Via Proxy

```bash
ab -n 10000 -c 100 http://127.0.0.1:1235/
```

## Other Notes:

### build containers
 1. go_base_docker/build_docker_base.sh
 1. run_docker.sh

### boot2docker
 - **boot2docker ip** - to get ip address
 - **boot2docker ssh** - ssh to boot2docker box

### example simple requests

```bash
curl -vvv http://127.0.0.1:1235 -H 'Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245e;'
```
 
### example large requests
 
```bash
curl -vvv http://127.0.0.1:1235 --data 'thisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryveryloc'
curl -vvv http://127.0.0.1:1235 --data 'thisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryveryloc'
curl 'http://127.0.0.1:1235' -H 'Accept-Encoding: gzip,deflate,sdch' -H 'Accept-Language: en-US,en;q=0.8,fa;q=0.6' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Referer: https://www.google.co.uk/' -H 'Cookie: dynsoftup=452b8f23-fa46-11e3-9eba-28cfe9158b63; __utmb=110886291.4.10.1403464391; __utmc=110886291; __utmz=110886291.1403464391.7.6.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' --compressed
curl 'http://127.0.0.1:1235' -H 'Accept-Encoding: gzip,deflate,sdch' -H 'Accept-Language: en-US,en;q=0.8,fa;q=0.6' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Referer: https://www.google.co.uk/' -H 'Cookie: __utmb=110886291.4.10.1403464391; __utmc=110886291; __utmz=110886291.1403464391.7.6.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' --compressed
````


### example requests for upgrades
 - **Request for gradual (Upgrade happens after 6 request if "percentageTransitionPerRequest" is set to 0.5)** 
```bash
curl -v 'http://127.0.0.1:1235/' -H 'Cookie:   transition=952c8557-2088-11e4-87e3-600308a8245e;'
```

 - **Request for Session Upgrade**
```bash
curl -v 'http://127.0.0.1:1235/' -H 'Cookie:   dynsoftup=9f0ef721-2228-11e4-9770-600308a8245e;'
```

### proxy to real server

i.e. http://nginx.com/blog/http-keepalives-and-web-performance/

1. Update proxy to call nginx.com:80
1. Run following curl command

```bash
curl -vvv http://127.0.0.1:1235/blog/http-keepalives-and-web-performance/ -H "Host: nginx.com"
```

### testing config REST services using Chrome DHC

The easiest way to test the config services is to use [DHC](https://chrome.google.com/webstore/detail/dhc-rest-http-api-client/aejoelaoggembcahagimdiliamlcdmfm) and import the example calls from DHC_Chrome_Extension_Config_Server_REST_Examples.json in the project root.


### testing config REST services using curl

 - **PUT Request for Concurrent upgrade**
```bash
curl http://127.0.0.1:9090/server -X PUT --data '{"cluster": {"servers":[{"ip": "127.0.0.1", "port": 1037},{"ip": "127.0.0.1", "port": 1038},{"ip": "127.0.0.1", "port": 1039}],"version": 1.1,"upgradeTransition": { "mode": "CONCURRENT" }}}'
````

 - **PUT Request using a config file**
```bash
curl -i -H "Accept: application/json" -X PUT -d @config/config_script.json localhost:9090/server
```

 - **GET Request for getting the list of all versions of the clusters**
```bash
curl http://127.0.0.1:9090/server/ -X GET 
```

 - **GET Request for getting a specific cluster from the list of clusters**
```bash
curl http://127.0.0.1:9090/server/e3011edf-2249-11e4-9b84-600308a8245e -X GET 
```

 - **DELETE Request for deleting a specific cluster from the list of clusters**
```bash
curl http://127.0.0.1:9090/server/e3011edf-2249-11e4-9b84-600308a8245e -X DELETE
```


### test running in boot2docker

To make boot2docker work it needs a copy of the config file read by the docker proxy container, as follows:

```bash
boot2docker ssh "mkdir /home/docker/config; cd /home/docker/config; cat << EOF > config.json
> {
>     "proxy": {
>         "ip": "localhost",
>         "port": 1234
>     },
>     "server_range":{
>         "ip": "127.0.0.1",
>         "port": 1024,
>         "clusterSize": "8"
>     }
> }"
```

```bash
# from curl
curl -vvv http://$(boot2docker ip 2>/dev/null):1234 -H 'Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245e;'

# from wrk
./wrk -t200 -c200 -d10 --latency http://$(boot2docker ip 2>/dev/null):1234
```

# Gradual Transition UUID to Percentages

- a37a290f-2088-11e4-b3a6-600308a8245e => 1
- a37a2633-2088-11e4-b3a6-600308a8245e => 2
- 952c8557-2088-11e4-87e3-600308a8245e => 3
- 10671a7f-2087-11e4-bf9e-600308a8245e => 4
- a37a25b4-2088-11e4-b3a6-600308a8245e => 5
- 6024b0c6-2089-11e4-a4ef-600308a8245e => 8
- 6024b0a9-2089-11e4-a4ef-600308a8245e => 9
- 73f9e08a-2089-11e4-b5af-600308a8245e => 10
- 73f9dfb0-2089-11e4-b5af-600308a8245e => 12
- 106719e4-2087-11e4-bf9e-600308a8245e => 80
- a37a2717-2088-11e4-b3a6-600308a8245e => 85
- a37a2826-2088-11e4-b3a6-600308a8245e => 90

# Test Coverage

export PATH=$PATH:$GOPATH/bin
export GOPATH=$GOPATH:`PWD`

go get github.com/axw/gocov/gocov
go get gopkg.in/matm/v1/gocov-html
gocov test -v proxy | gocov-html > coverage.html

# Run Docker with Go as Terminal

1. vagrant up
1. vagrant ssh
1. docker run -i -v /vagrant:/vagrant -t samirabloom/docker-go /bin/bash
1. cd /vagrant

# Debug Docker Containers

1. commit container create as follows: docker commit <container name / id> <tag>
1. run container in interactive mode as follows: docker run -i -t <tag> /bin/bash

# Curl Different Couchbase Versions

# curl 2.5.1
curl -v http://Administrator:password@192.168.50.50:8091/pools/default/buckets/beer-sample

# curl 3.0.0
curl -v http://Administrator:password@192.168.50.60:8091/pools/default/buckets/beer-sample
