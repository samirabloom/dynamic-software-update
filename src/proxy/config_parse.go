package proxy

import (
	"code.google.com/p/go-uuid/uuid"
	"io/ioutil"
	"net"
	"encoding/json"
	"fmt"
	"errors"
	"proxy/log"
	"proxy/stages"
)

// ==== PARSE CONFIG - START

func loadConfig(configFile string) (*Proxy, error) {
	return parseConfigFile(readConfigFile(configFile), parseProxy, parseClusters(func() uuid.UUID { return uuid.NewUUID() }))
}

func readConfigFile(configFile string) []byte {
	jsonConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.LoggerFactory().Error("Error %s reading config file [%s]", err, configFile)
	}
	return jsonConfig
}

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseClusters func(map[string]interface{}) (*stages.Clusters, error)) (proxy *Proxy, err error) {
	// parse json object
	var jsonConfig = make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		log.LoggerFactory().Error("Error %s parsing config file:\n%s", err.Error(), jsonData)
	}

	tcpProxyLocalAddress, proxyParseErr := parseProxy(jsonConfig)
	if proxyParseErr == nil {
		configServicePort, parseConfigServiceErr := parseConfigService(jsonConfig)
		if parseConfigServiceErr == nil {
			clusters, clusterParseErr := parseClusters(jsonConfig)
			if clusterParseErr == nil {
				// create load balancer
				proxy = &Proxy{
					frontendAddr: tcpProxyLocalAddress,
					configServicePort: configServicePort,
					clusters: clusters,
					stop: make(chan bool),
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
			log.LoggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}
	} else {
		errorMessage := "Invalid proxy configuration - \"proxy\" config missing"
		log.LoggerFactory().Error(errorMessage)
		err = errors.New(errorMessage)
	}

	return tcpProxyLocalAddress, err
}

func parseConfigService(jsonConfig map[string]interface{}) (int, error) {
	var (
		err error
		configServicePort int
	)

	if jsonConfig["configService"] != nil {
		var configServiceConfig map[string]interface{} = jsonConfig["configService"].(map[string]interface{})
		if configServiceConfig["port"] != nil {
			configServicePort = int(configServiceConfig["port"].(float64))
		} else {
			errorMessage := "Invalid config service configuration - \"port\" is missing from \"configService\" config"
			log.LoggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}
	} else {
		errorMessage := "Invalid proxy configuration - \"configService\" config missing"
		log.LoggerFactory().Error(errorMessage)
		err = errors.New(errorMessage)
	}
	return configServicePort, err
}

func parseClusters(uuidGenerator func() uuid.UUID) func(map[string]interface{}) (*stages.Clusters, error) {
	return func(jsonConfig map[string]interface{}) (*stages.Clusters, error) {
		var (
			err error
			router *stages.Cluster
			clusters *stages.Clusters
		)

		clusterConfiguration := jsonConfig["cluster"]
		if clusterConfiguration != nil {
			router, err = parseCluster(uuidGenerator)(clusterConfiguration.(map[string]interface{}))
			if err == nil {
				clusters = &stages.Clusters{}
				clusters.Add(router)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"cluster\" config missing"
			log.LoggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}

		return clusters, err
	}
}

func parseCluster(uuidGenerator func() uuid.UUID) func(map[string]interface{}) (*stages.Cluster, error) {
	return func(clusterConfiguration map[string]interface{}) (*stages.Cluster, error) {
		var (
			err error
			backendAddresses []*net.TCPAddr
			version float64
			sessionTimeout int64
			mode stages.TransitionMode
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
						log.LoggerFactory().Error(errorMessage)
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

				upgradeTransitionConfig := clusterConfiguration["upgradeTransition"]
				if upgradeTransitionConfig != nil {
					upgradeTransition := upgradeTransitionConfig.(map[string]interface{})

					modeConfig := upgradeTransition["mode"]
					if modeConfig != nil {
						mode = stages.ModesCodeToMode[modeConfig.(string)]
					} else {
						mode = stages.SessionMode
					}

					if mode != 0 {
						sessionTimeoutConfig := upgradeTransition["sessionTimeout"]
						if mode == stages.SessionMode {
							if sessionTimeoutConfig != nil {
								sessionTimeout = int64(sessionTimeoutConfig.(float64))
							} else {
								errorMessage := "Invalid cluster configuration - \"sessionTimeout\" is missing from \"upgradeTransition\" config"
								log.LoggerFactory().Error(errorMessage)
								err = errors.New(errorMessage)
							}
						} else if sessionTimeoutConfig != nil {
							errorMessage := "Invalid cluster configuration - \"sessionTimeout\" should not be specified when \"mode\" is \"INSTANT\""
							log.LoggerFactory().Error(errorMessage)
							err = errors.New(errorMessage)
						}
					} else {
						errorMessage := "Invalid cluster configuration - \"upgradeTransition.mode\" should be \"" + stages.ModesModeToCode[stages.InstantMode] + "\" or \"" + stages.ModesModeToCode[stages.SessionMode] + "\""
						log.LoggerFactory().Error(errorMessage)
						err = errors.New(errorMessage)
					}
				} else {
					sessionTimeout = 0
					mode = stages.InstantMode
				}
			} else {
				errorMessage := "Invalid cluster configuration - \"servers\" list must contain at least one entry"
				log.LoggerFactory().Error(errorMessage)
				err = errors.New(errorMessage)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"servers\" list missing from \"cluster\" config"
			log.LoggerFactory().Error(errorMessage)
			err = errors.New(errorMessage)
		}

		return &stages.Cluster{BackendAddresses: backendAddresses, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: sessionTimeout, Mode: mode, Version: version}, err
	}
}

func serialiseCluster(cluster *stages.Cluster) map[string]interface{} {
	jsonConfig := map[string]interface{}{}

	if cluster != nil {
		var serversConfig []interface{} = make([]interface{}, len(cluster.BackendAddresses))
		for index, backendAddress := range cluster.BackendAddresses {
			serversConfig[index] = map[string]interface{}{"ip": backendAddress.IP, "port": backendAddress.Port}
		}
		jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"uuid": cluster.Uuid.String(), "servers": serversConfig, "version": cluster.Version, "upgradeTransition": map[string]interface{}{"sessionTimeout": cluster.SessionTimeout, "mode": stages.ModesModeToCode[cluster.Mode]}}}
	}

	return jsonConfig
}

// ==== PARSE CONFIG - END
