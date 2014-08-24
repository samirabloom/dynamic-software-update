package proxy

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"proxy/contexts"
	"proxy/log"
)

// ==== PARSE CONFIG - START

func loadConfig(configFile string) (*Proxy, error) {
	jsonData, err := readConfigFile(configFile)
	if err == nil {
		return parseConfigFile(jsonData, parseProxy, parseConfigService, parseClusters(func() uuid.UUID { return uuid.NewUUID() }, true))
	} else {
		return nil, err
	}
}

func readConfigFile(configFile string) ([]byte, error) {
	jsonConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error %s reading config file [%s]", err, configFile))
	}
	return jsonConfig, nil
}

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseConfigService func(map[string]interface{}) (int, error), parseClusters func(map[string]interface{}) (*contexts.Clusters, error)) (proxy *Proxy, err error) {
	// parse json object
	var jsonConfig = make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error %s parsing config file:\n%s", err.Error(), jsonData))
	} else {
		tcpProxyLocalAddress, proxyParseErr := parseProxy(jsonConfig)
		if proxyParseErr == nil {
			configServicePort, parseConfigServiceErr := parseConfigService(jsonConfig)
			if parseConfigServiceErr == nil {
				clusters, clusterParseErr := parseClusters(jsonConfig)
				if clusterParseErr == nil {
					// create load balancer
					proxy = &Proxy{
						frontendAddr:      tcpProxyLocalAddress,
						configServicePort: configServicePort,
						clusters:          clusters,
						stop:              make(chan bool),
					}
					log.LoggerFactory().Notice("Parsed config file:\n%s\nas:\n%s", jsonData, proxy)

					return proxy, nil
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
}

func parseProxy(jsonConfig map[string]interface{}) (*net.TCPAddr, error) {
	var (
		err                  error
		tcpProxyLocalAddress *net.TCPAddr
	)

	if jsonConfig["proxy"] != nil {
		var proxyConfig map[string]interface{} = jsonConfig["proxy"].(map[string]interface{})
		tcpProxyLocalAddress, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", proxyConfig["port"]))
		if err != nil {
			errorMessage := "Invalid proxy address [" + fmt.Sprintf(":%v", proxyConfig["port"]) + "] - " + err.Error()
			err = errors.New(errorMessage)
		}
	} else {
		errorMessage := "Invalid proxy configuration - \"proxy\" config missing"
		err = errors.New(errorMessage)
	}

	return tcpProxyLocalAddress, err
}

func parseConfigService(jsonConfig map[string]interface{}) (int, error) {
	var (
		err               error
		configServicePort int
	)

	if jsonConfig["configService"] != nil {
		var configServiceConfig map[string]interface{} = jsonConfig["configService"].(map[string]interface{})
		if configServiceConfig["port"] != nil {
			configServicePort = int(configServiceConfig["port"].(float64))
		} else {
			errorMessage := "Invalid config service configuration - \"port\" is missing from \"configService\" config"
			err = errors.New(errorMessage)
		}
	} else {
		errorMessage := "Invalid proxy configuration - \"configService\" config missing"
		err = errors.New(errorMessage)
	}
	return configServicePort, err
}

func parseClusters(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}) (*contexts.Clusters, error) {
	return func(jsonConfig map[string]interface{}) (*contexts.Clusters, error) {
		var (
			err      error
			router   *contexts.Cluster
			clusters *contexts.Clusters
		)

		clusterConfiguration := jsonConfig["cluster"]
		if clusterConfiguration != nil {
			router, err = parseCluster(uuidGenerator, initialCluster)(clusterConfiguration.(map[string]interface{}))
			if err == nil {
				clusters = &contexts.Clusters{}
				clusters.Add(router)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"cluster\" config missing"
			err = errors.New(errorMessage)
		}

		return clusters, err
	}
}

func parseCluster(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}) (*contexts.Cluster, error) {
	return func(clusterConfiguration map[string]interface{}) (*contexts.Cluster, error) {
		var (
			err                            error
			connection 					   *net.TCPAddr
			backendAddresses               []*contexts.BackendAddress
			version                        float64
			sessionTimeout                 int64
			percentageTransitionPerRequest float64
			mode                           contexts.TransitionMode
			uuidValue                      uuid.UUID
		)

		serversConfiguration := clusterConfiguration["servers"]
		if serversConfiguration != nil {
			servers := serversConfiguration.([]interface{})
			if len(servers) > 0 {
				backendAddresses = make([]*contexts.BackendAddress, len(servers))
				for index := range servers {
					var server map[string]interface{} = servers[index].(map[string]interface{})
					connection, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", server["hostname"], server["port"]))
					if err != nil {
						errorMessage := "Invalid server address [" + fmt.Sprintf("%s:%v", server["hostname"], server["port"]) + "] - " + err.Error()
						err = errors.New(errorMessage)
					} else {
						backendAddresses[index] = &contexts.BackendAddress{Address: connection, Host: fmt.Sprintf("%s", server["hostname"]), Port: fmt.Sprintf("%v", server["port"])}
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

				upgradeTransitionConfig := clusterConfiguration["upgradeTransition"]
				if upgradeTransitionConfig != nil {
					if initialCluster {
						errorMessage := "Invalid cluster configuration - \"upgradeTransition\" can not be specified for the intial cluster"
						err = errors.New(errorMessage)
					} else {
						upgradeTransition := upgradeTransitionConfig.(map[string]interface{})

						modeConfig := upgradeTransition["mode"]
						if modeConfig != nil {
							mode = contexts.ModesCodeToMode[modeConfig.(string)]
						} else {
							mode = contexts.SessionMode
						}

						if mode != 0 {
							sessionTimeoutConfig := upgradeTransition["sessionTimeout"]
							if mode == contexts.SessionMode {
								if sessionTimeoutConfig != nil {
									sessionTimeout = int64(sessionTimeoutConfig.(float64))
								} else {
									errorMessage := "Invalid cluster configuration - \"sessionTimeout\" must be specified in \"upgradeTransition\" for mode \"SESSION\""
									err = errors.New(errorMessage)
								}
							} else if sessionTimeoutConfig != nil {
								errorMessage := "Invalid cluster configuration - \"sessionTimeout\" should not be specified when \"mode\" is not \"SESSION\""
								err = errors.New(errorMessage)
							}

							percentageTransitionPerRequestConfig := upgradeTransition["percentageTransitionPerRequest"]
							if mode == contexts.GradualMode {
								if percentageTransitionPerRequestConfig != nil {
									percentageTransitionPerRequest = percentageTransitionPerRequestConfig.(float64)
								} else {
									errorMessage := "Invalid cluster configuration - \"percentageTransitionPerRequest\" must be specified in \"upgradeTransition\" for mode \"GRADUAL\""
									err = errors.New(errorMessage)
								}
							} else if percentageTransitionPerRequestConfig != nil {
								errorMessage := "Invalid cluster configuration - \"percentageTransitionPerRequest\" should not be specified when \"mode\" is not \"GRADUAL\""
								err = errors.New(errorMessage)
							}
						} else {
							errorMessage := "Invalid cluster configuration - \"upgradeTransition.mode\" should be \"INSTANT\", \"SESSION\", \"GRADUAL\" or \"CONCURRENT\""
							err = errors.New(errorMessage)
						}
					}
				} else {
					sessionTimeout = 0
					percentageTransitionPerRequest = 0
					mode = contexts.InstantMode
				}
			} else {
				errorMessage := "Invalid cluster configuration - \"servers\" list must contain at least one entry"
				err = errors.New(errorMessage)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"servers\" list missing from \"cluster\" config"
			err = errors.New(errorMessage)
		}

		return &contexts.Cluster{BackendAddresses: backendAddresses, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: sessionTimeout, PercentageTransitionPerRequest: percentageTransitionPerRequest, Mode: mode, Version: version}, err
	}
}

func serialiseCluster(cluster *contexts.Cluster) map[string]interface{} {
	jsonConfig := map[string]interface{}{}

	if cluster != nil {
		var serversConfig []interface{} = make([]interface{}, len(cluster.BackendAddresses))
		for index, backendAddress := range cluster.BackendAddresses {
			serversConfig[index] = map[string]interface{}{"hostname": backendAddress.Host, "port": backendAddress.Address.Port}
		}
		upgradeTransition := map[string]interface{}{"mode": contexts.ModesModeToCode[cluster.Mode]}
		switch cluster.Mode {
		case contexts.SessionMode:
		{
			upgradeTransition = map[string]interface{}{"mode": contexts.ModesModeToCode[cluster.Mode], "sessionTimeout": cluster.SessionTimeout}
		}
		case contexts.GradualMode:
		{
			upgradeTransition = map[string]interface{}{"mode": contexts.ModesModeToCode[cluster.Mode], "percentageTransitionPerRequest": cluster.PercentageTransitionPerRequest}
		}
		}
		jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"uuid": cluster.Uuid.String(), "servers": serversConfig, "version": cluster.Version, "upgradeTransition": upgradeTransition}}
	}

	return jsonConfig
}

// ==== PARSE CONFIG - END
