dynamic-software-update
=======================

## Build & Run

```bash
./build_run.sh
```

This will run the proxy on port 1234 and will run a server on port 1024.  The server has a chunked HTTP response with a 230ms forced delay.
 
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
./wrk -t200 -c200 -d10 --latency http://127.0.0.1:1234
```

## ApacheBench
 
### 1. Test Server Directly
 
```bash
ab -n 10000 -c 100 http://127.0.0.1:1024/
```

### 2. Test Server Via Proxy

```bash
ab -n 10000 -c 100 http://127.0.0.1:1234/
```

## Other Notes:

### build containers
 1. go_base_docker/build_docker_base.sh
 1. run_docker.sh

### boot2docker
 - **boot2docker ip** - to get ip address
 - **boot2docker ssh** - ssh to boot2docker box

### example simple requests

curl -vvv http://127.0.0.1:1234 -H 'Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245e;'
 
### example large requests
 
```bash
curl -vvv http://127.0.0.1:1234/JVMInternals.html -H "Host: blog.jamesdbloom.com:80"
curl -vvv http://127.0.0.1:1234 --data 'thisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryveryloc'
curl -vvv http://127.0.0.1:1234 --data 'thisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylocthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryverylongparameterthatgoesonandonandonthisisaveryveryveryveryveryveryveryveryveryloc'
curl 'http://127.0.0.1:1234' -H 'Accept-Encoding: gzip,deflate,sdch' -H 'Accept-Language: en-US,en;q=0.8,fa;q=0.6' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Referer: https://www.google.co.uk/' -H 'Cookie: dynsoftup=452b8f23-fa46-11e3-9eba-28cfe9158b63; __utmb=110886291.4.10.1403464391; __utmc=110886291; __utmz=110886291.1403464391.7.6.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' --compressed
curl 'http://127.0.0.1:1234' -H 'Accept-Encoding: gzip,deflate,sdch' -H 'Accept-Language: en-US,en;q=0.8,fa;q=0.6' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Referer: https://www.google.co.uk/' -H 'Cookie: __utmb=110886291.4.10.1403464391; __utmc=110886291; __utmz=110886291.1403464391.7.6.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' --compressed
````