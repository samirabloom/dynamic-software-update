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
                                                       {"hostname": "127.0.0.1", "port": 1034},
                                                       {"hostname": "127.0.0.1", "port": 1035}
                                                   ],
                                                   "version": "1.0"
                                               }
                                           }

  -logLevel="WARN": Set the log level as "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO" or "DEBUG"
  
  -h: Display this message
```

For example:

```bash
proxy -logLevel=INFO -configFile="config/config_script.json"
```

