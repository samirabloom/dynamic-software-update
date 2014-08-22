The proxy provides a simple REST API to support dynamically updating the cluster configuration as follows:
 
* **PUT /configuration/cluster** - adds a new cluster configuration
* **GET /configuration/cluster/{clusterId}** - gets a single cluster configuration
* **GET /configuration/cluster** - gets all cluster configurations
* **DELETE /configuration/cluster/{clusterId}** - deletes a single cluster configuration

## HTTP Response Codes

* **202 Accepted** - a new cluster entity is successfully added or deleted 
* **200 OK** - cluster(s) entity is successfully returned   
* **404 Not Found** - cluster id is invalid
* **400 Bad Request** - request syntax is  invalid

## PUT - /configuration/cluster
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

##### cluster.version
Type: `Number` Default value: `0`

This value specifies the cluster version. If no version is specified, the version defaults to `0`. 

##### cluster.upgradeTransition
Type: `Object` Default value: `{ mode: "INSTANT" }`

This value allows the configuration of the upgrade transition. If no `upgradeTransition` is specified, the upgrade transition mode defaults to `INSTANT`.

##### cluster.upgradeTransition.mode
Type: `String` Default value: `SESSION`

This value specifies the upgrade transition mode and support the following values: `INSTANT`, `SESSION`, `GRADUAL`, `CONCURRENT`. If no `upgradeTransition mode` is specified, the mode defaults to `SESSION`.

##### cluster.upgradeTransition.sessionTimeout
Type: `Number` Default value: `undefined`

This value specifies the timeout period assigned to the `SESSION` transition mode.

##### cluster.upgradeTransition.percentageTransitionPerRequest
Type: `Number` Default value: `undefined`

This value specifies the transition percentage associated with each request in the `GRADUAL` transition mode. 
 
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
## GET - /configuration/cluster/{clusterId}

To get a single cluster configuration make a GET request to `/configuration/cluster/{clusterId}`. 

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
    "upgradeTransition": {
      "mode": ""
      "sessionTimeout": 0  // only returned when `mode` is "SESSION" 
      "percentageTransitionPerRequest": 0  // only returned when `mode` is "GRADUAL"
    },
    "uuid": "",
    "version": 0
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

## GET - /configuration/cluster

To get all the cluster configurations make a GET request with no cluster id `/configuration/cluster/`.

### Response Body

````js

[
  {
    "cluster": {
      "servers": [
        {
          "ip": "",
          "port": 0
        }
       ],
      "upgradeTransition": {
        "mode": ""
        "sessionTimeout": 0  // only returned when `mode` is "SESSION" 
        "percentageTransitionPerRequest": 0  // only returned when `mode` is "GRADUAL"
      },
      "uuid": "",
      "version": 0
    }
  },
  {
    "cluster": {
      "servers": [
        {
          "ip": "",
          "port": 0
        },
        {
          "ip": "",
          "port": 0
        }
       ],
      "upgradeTransition": {
        "mode": "CONCURRENT"
      },
      "uuid": "",
      "version": 0
    }
  }
]
```

### Example

##### Request

For example the following `curl` request would get a list of all cluster configurations

```bash
curl http://127.0.0.1:9090/configuration/cluster/ -X GET
```
##### Response

```js
[
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
      "uuid": "1f6a0854-2608-11e4-ab79-600308a8245e", 
      "version": 1.1
    }
  }, 
  {
    "cluster": {
      "servers": [
        {
          "ip": "127.0.0.1", 
          "port": 1037
        }, 
        {
          "ip": "127.0.0.1", 
          "port": 1039
        }
      ], 
      "upgradeTransition": {
        "mode": "CONCURRENT"
      }, 
      "uuid": "01386f1f-2608-11e4-ab79-600308a8245e", 
      "version": 1.1
    }
  }, 
  {
    "cluster": {
      "servers": [
        {
          "ip": "127.0.0.1", 
          "port": 1034
        }, 
        {
          "ip": "127.0.0.1", 
          "port": 1035
        }
      ], 
      "upgradeTransition": {
        "mode": "INSTANT"
      }, 
      "uuid": "ffde36ce-2607-11e4-ab79-600308a8245e", 
      "version": 1
    }
  }
]
```

For example the response when using curl is as follows:

```bash
HTTP/1.1 200 OK
Date: Sun, 17 Aug 2014 12:28:55 GMT
Content-Length: 583
Content-Type: text/plain; charset=utf-8
 
[{"cluster":{"servers":[{"ip":"127.0.0.1","port":1036},{"ip":"127.0.0.1","port":1038}],"upgradeTransition":{"mode":"SESSION","sessionTimeout":60},"uuid":"1f6a0854-2608-11e4-ab79-600308a8245e","version":1.1}},{"cluster":{"servers":[{"ip":"127.0.0.1","port":1037},{"ip":"127.0.0.1","port":1039}],"upgradeTransition":{"mode":"CONCURRENT"},"uuid":"01386f1f-2608-11e4-ab79-600308a8245e","version":1.1}},{"cluster":{"servers":[{"ip":"127.0.0.1","port":1034},{"ip":"127.0.0.1","port":1035}],"upgradeTransition":{"mode":"INSTANT"},"uuid":"ffde36ce-2607-11e4-ab79-600308a8245e","version":1}}]
```

## DELETE - /configuration/cluster/{clusterId}

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
