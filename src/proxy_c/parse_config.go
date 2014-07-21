package proxy_c

import (
	"code.google.com/p/go-uuid/uuid"
	"io/ioutil"
	"net"
	"encoding/json"
	"fmt"
	"errors"
)

// ==== PARSE CONFIG - START

func loadConfig(configFile *string) (*LoadBalancer, error) {
	return parseConfigFile(readConfigFile(configFile), parseProxy, parseRoutingContexts(func() uuid.UUID {
			return uuid.NewUUID()
		}))
}

func readConfigFile(configFile *string) []byte {
	jsonConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		loggerFactory().Error("Error %s reading config file [%s]", err, *configFile)
	}
	return jsonConfig
}

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseRoutingContexts func(map[string]interface{}) (*RoutingContexts, error)) (loadBalancer *LoadBalancer, err error) {
	// parse json object
	var jsonConfig = make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		loggerFactory().Error("Error %s parsing config file:\n%s", err.Error(), jsonData)
	}

	tcpProxyLocalAddress, proxyParseErr := parseProxy(jsonConfig)
	if proxyParseErr == nil {
		configServicePort, parseConfigServiceErr := parseConfigService(jsonConfig)
		if parseConfigServiceErr == nil {
			routingContexts, clusterParseErr := parseRoutingContexts(jsonConfig)
			if clusterParseErr == nil {
				// create load balancer
				loadBalancer = &LoadBalancer{
					frontendAddr: tcpProxyLocalAddress,
					configServicePort: configServicePort,
					routingContexts: routingContexts,
					stop: make(chan bool),
				}
				loggerFactory().Info("Parsed config file:\n%s\nas:\n%s", jsonData, loadBalancer)

				return loadBalancer, nil
			} else {
				return nil, clusterParseErr
			}
		} else {
			return nil, parseConfigServiceErr
		}
	} else {
		return nil, proxyParseErr
	}
}

func parseProxy(jsonConfig map[string]interface{}) (*net.TCPAddr, error) {
	var (
		err error
		tcpProxyLocalAddress *net.TCPAddr
	)

	if jsonConfig["proxy"] != nil {
		var proxyConfig map[string]interface{} = jsonConfig["proxy"].(map[string]interface{})
		tcpProxyLocalAddress, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]))
		if err != nil {
			errorMessage := "Invalid proxy address [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "] - " + err.Error()
			loggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}
	} else {
		errorMessage := "Invalid proxy configuration - \"proxy\" JSON field missing or invalid"
		loggerFactory().Error(errorMessage)
		err = errors.New(errorMessage)
	}

	return tcpProxyLocalAddress, err
}

func parseConfigService(jsonConfig map[string]interface{}) (float64, error) {
	var (
		err error
		configServicePort float64
	)

	if jsonConfig["configService"] != nil {
		var configServiceConfig map[string]interface{} = jsonConfig["configService"].(map[string]interface{})
		configServicePort = configServiceConfig["port"].(float64)
		if configServiceConfig["port"] == nil {
			errorMessage := "Invalid configService port ["+fmt.Sprintf("%v", configServiceConfig["port"])+"]"
			loggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}
	}
	return configServicePort, err
}

func parseRoutingContexts(uuidGenerator func() uuid.UUID) func(map[string]interface{}) (*RoutingContexts, error) {
	return func(jsonConfig map[string]interface{}) (*RoutingContexts, error) {
		var (
			err error
			router *RoutingContext
			routingContexts *RoutingContexts
		)

		clusterConfiguration := jsonConfig["cluster"]
		if clusterConfiguration != nil {
			router, err = parseRoutingContext(uuidGenerator)(clusterConfiguration.(map[string]interface{}))
			if err == nil {
				routingContexts = &RoutingContexts{}
				routingContexts.Add(router)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"cluster\" config missing"
			loggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}

		return routingContexts, err
	}
}

func parseRoutingContext(uuidGenerator func() uuid.UUID) func(map[string]interface{}) (*RoutingContext, error) {
	return func(clusterConfiguration map[string]interface{}) (*RoutingContext, error) {
		var (
			err error
			backendAddresses []*net.TCPAddr
			version float64
			uuidValue uuid.UUID
		)

		serversConfiguration := clusterConfiguration["servers"]
		if serversConfiguration != nil {
			servers := serversConfiguration.([]interface{})
			if len(servers) > 0 {
				backendAddresses = make([]*net.TCPAddr, len(servers))
				for index := range servers {
					var server map[string]interface{} = servers[index].(map[string]interface{})
					backendAddresses[index], err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", server["ip"], server["port"]))
					if err != nil {
						errorMessage := "Invalid server address [" + fmt.Sprintf("%s:%v", server["ip"], server["port"]) + "] - " + err.Error()
						loggerFactory().Error(errorMessage)
						err = errors.New(errorMessage)
					}
				}
				uuidConfig := clusterConfiguration["uuid"]
				if uuidConfig != nil {
					uuidValue = uuid.Parse(uuidConfig.(string))
				} else {
					uuidValue = uuidGenerator()
				}
				versionConfig := clusterConfiguration["version"]
				if versionConfig != nil {
					version = versionConfig.(float64)
				} else {
					version = 0.0
				}
			} else {
				errorMessage := "Invalid cluster configuration - \"servers\" list must contain at least one entry"
				loggerFactory().Error(errorMessage)
				err = errors.New(errorMessage)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"servers\" list missing from \"cluster\" config"
			loggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}

		return &RoutingContext{backendAddresses: backendAddresses, requestCounter: -1, uuid: uuidValue, version: version}, err
	}
}

func serialiseRoutingContext(routingContext *RoutingContext) map[string]interface{} {
	jsonConfig := map[string]interface{}{}

	if routingContext != nil {
		var serversConfig []interface{} = make([]interface{}, len(routingContext.backendAddresses))
		for index, backendAddress := range routingContext.backendAddresses {
			serversConfig[index] = map[string]interface{}{"ip": backendAddress.IP, "port": backendAddress.Port}
		}
		jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"uuid": routingContext.uuid.String(), "servers": serversConfig, "version": routingContext.version}}
	}

	return jsonConfig
}

// ==== PARSE CONFIG - END
