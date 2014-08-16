The proxy provides a simple REST API to support dynamically updating the cluster configuration as follows:
 
* **PUT /configuration/cluster** - adds a new cluster configuration
* **GET /configuration/cluster/{clusterId}** - gets a single cluster configuration
* **GET /configuration/cluster** - gets all cluster configurations
* **DELETE /configuration/cluster/{clusterId}** - deletes a single cluster configuration

## HTTP Response Codes

* **202 StatusAccepted** - a new cluster entity is successfully added or deleted 
* **200 StatusOK** - cluster(s) entity is successfully returned   
* **404 StatusNotFound** - cluster id is invalid
* **400 StatusBadRequest** - request syntax is  invalid

## /configuration/cluster

This endpoint represents the cluster entity and supports the following HTTP methods:

## PUT
To add a new cluster make a PUT request to `/configuration/cluster`.

###Request Body
 
```js
{
  "cluster": {
    "servers": [
      {
        "ip": "", 
        "port": 0
      }
    ], 
    "version": 0, 
    "upgradeTransition": {
      "mode": ""  // allowed values are "INSTANT", "SESSION", "GRADUAL", "CONCURRENT"
      "sessionTimeout": 0  // only supported for a `mode` value of "SESSION" 
      "percentageTransitionPerRequest": 0  // only supported for a `mode` value of "GRADUAL"
    }
  }
}
```

##### cluster.servers

Type: `Array ` Default value: `[]`

This value specifies the list of servers in the cluster

##### cluster.servers[i].ip
Type: `String` Default value: `undefined`

This value specifies the ip address or hostname of a server in the cluster

##### cluster.servers[i].port
Type: `Number` Default value: `undefined`

This value specifies the port of a server in the cluster

##### cluster.upgradeTransition.mode
Type: `String` Default value: `INSTANT`

This value specifies the upgrade transition mode and support the following values: `INSTANT`, `SESSION`, `GRADUAL`, `CONCURRENT`

##### cluster.upgradeTransition.sessionTimeout
Type: `String` Default value: `undefined`

This value specifies the timeout period...
 

###Response Body

A cluster id is returned representing the new cluster entity that has been added. 

###Example

##### Request

For example the following JSON would set up a new cluster with two `servers` and `SESSION` upgrade transition:

```js
{
  "cluster": {
    "servers": [
      {
        "ip": "127.0.0.1", 
        "port": 1036
      },  
      {
        "ip": "127.0.0.1", 
        "port": 1038
      }
    ], 
    "version": 1.1, 
    "upgradeTransition": {
      "mode": "SESSION", 
      "sessionTimeout": 60
    }
  }
}
```

To send this request with `Curl` use the following syntax:

```bash
curl http://127.0.0.1:9090/configuration/cluster -X PUT --data '{"cluster": {"servers":[{"ip": "127.0.0.1", "port": 1036},{"ip": "127.0.0.1", "port": 1038}],"version": 1.1,"upgradeTransition": { "mode": "SESSION", "sessionTimeout": 60 }}}'
```

##### Response

```bash
HTTP/1.1 202 Accepted
Date: Sat, 16 Aug 2014 19:54:21 GMT
Content-Length: 36
Content-Type: text/plain; charset=utf-8
 
1dcbb083-257f-11e4-bcbc-600308a8245e
```
## GET

To get a single cluster configuration make a GET request to `/configuration/cluster/{clusterId}`. To get all the cluster configurations make a GET request with no cluster id `/configuration/cluster/`.

### Response Body

```js
{
  "cluster": {
    "servers": [
      {
        "ip": "", 
        "port": 0
      }
    ], 
    "version": 0, 
    "upgradeTransition": {
      "mode": ""  
      "sessionTimeout": 0  // only returned when `mode` is "SESSION" 
      "percentageTransitionPerRequest": 0  // only returned when `mode` is "GRADUAL"
    }
  }
}
```

### Example

##### Request

For example the following `curl` request would get the cluster configuration with cluster id `1dcbb083-257f-11e4-bcbc-600308a8245e`:

```bash
curl http://127.0.0.1:9090/configuration/cluster/1dcbb083-257f-11e4-bcbc-600308a8245e -X GET
```
##### Response

```js
{
  "cluster": {
    "servers": [
      {
        "ip": "127.0.0.1", 
        "port": 1036
      }, 
      {
        "ip": "127.0.0.1", 
        "port": 1038
      }
    ], 
    "upgradeTransition": {
      "mode": "SESSION", 
      "sessionTimeout": 60
    }, 
    "uuid": "016ca2cd-2585-11e4-ab5c-600308a8245e", 
    "version": 1.1
  }
}
```

For example the response when using curl is as follows:

```bash
HTTP/1.1 200 OK
Date: Sat, 16 Aug 2014 20:37:42 GMT
Content-Length: 206
Content-Type: text/plain; charset=utf-8

{"cluster":{"servers":[{"ip":"127.0.0.1","port":1036},{"ip":"127.0.0.1","port":1038}],"upgradeTransition":{"mode":"SESSION","sessionTimeout":60},"uuid":"016ca2cd-2585-11e4-ab5c-600308a8245e","version":1.1}}
```

## DELETE

To delete a single cluster configuration make a DELETE request to `/configuration/cluster/{clusterId}`.

### Example

##### Request

For example the following `curl` request would delete the cluster configuration with id `1dcbb083-257f-11e4-bcbc-600308a8245e`:

```bash
curl http://127.0.0.1:9090/configuration/cluster/1dcbb083-257f-11e4-bcbc-600308a8245e -X DELETE
```
##### Response

For example the response when using curl is as follows:

```bash
HTTP/1.1 202 Accepted
Date: Sat, 16 Aug 2014 21:28:38 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8

```



