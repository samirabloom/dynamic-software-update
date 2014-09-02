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
	"proxy/docker_client"
	"bytes"
	"io"
)

// ==== PARSE CONFIG - START

func LoadConfig(configFile string, outputStream io.Writer) (*Proxy, error) {
	jsonData, err := readConfigFile(configFile)
	if err == nil {
		return parseConfigFile(jsonData, parseProxy, parseConfigService, parseDockerHost, parseClusters(func() uuid.UUID { return uuid.NewUUID() }, true), outputStream)
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

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseConfigService func(map[string]interface{}) (int, error), parseDockerHost func(map[string]interface{}) (*docker_client.DockerHost, error), parseClusters func(map[string]interface{}, *docker_client.DockerHost, io.Writer) (*contexts.Clusters, error), outputStream io.Writer) (proxy *Proxy, err error) {
	var devNull bytes.Buffer

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
				dockerHost, parseDockerHostErr := parseDockerHost(jsonConfig)
				if parseDockerHostErr == nil {
					if dockerHost != nil && !dockerHost.Log {
						outputStream = &devNull
					}
					clusters, clusterParseErr := parseClusters(jsonConfig, dockerHost, outputStream)
					if clusterParseErr == nil {
						// create load balancer
						proxy = &Proxy{
							frontendAddr:      tcpProxyLocalAddress,
							configServicePort: configServicePort,
							dockerHost:        dockerHost,
							clusters:          clusters,
							stop:              make(chan bool),
						}
						log.LoggerFactory().Notice("Parsed config file:\n%s\nas:\n%s", jsonData, proxy)

						return proxy, nil
					} else {
						return nil, clusterParseErr
					}
				} else {
					return nil, parseDockerHostErr
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

func parseDockerHost(jsonConfig map[string]interface{}) (*docker_client.DockerHost, error) {
	var (
		err             error
		dockerHostIp    string
		dockerHostPort  int
		dockerHostLog   bool = true
		dockerHost      *docker_client.DockerHost
	)

	if jsonConfig["dockerHost"] != nil {
		var dockerHostConfig map[string]interface{} = jsonConfig["dockerHost"].(map[string]interface{})
		if dockerHostConfig["ip"] != nil {
			dockerHostIp = dockerHostConfig["ip"].(string)
		} else {
			errorMessage := "Invalid docker host configuration - \"ip\" is missing from \"dockerHost\" config"
			err = errors.New(errorMessage)
		}
		if dockerHostConfig["log"] != nil {
			dockerHostLog = dockerHostConfig["log"].(bool)
		}
		if err == nil {
			if dockerHostConfig["port"] != nil {
				dockerHostPort = int(dockerHostConfig["port"].(float64))
			} else {
				dockerHostPort = 2375
			}
			dockerHost = &docker_client.DockerHost{Ip: dockerHostIp, Port: dockerHostPort, Log: dockerHostLog}
		}
	}
	return dockerHost, err
}

func parseClusters(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}, *docker_client.DockerHost, io.Writer) (*contexts.Clusters, error) {
	return func(jsonConfig map[string]interface{}, dockerHost *docker_client.DockerHost, outputStream io.Writer) (*contexts.Clusters, error) {
		var (
			err      error
			router   *contexts.Cluster
			clusters *contexts.Clusters
		)

		clusterConfiguration := jsonConfig["cluster"]
		if clusterConfiguration != nil {
			router, err = parseCluster(uuidGenerator, initialCluster)(clusterConfiguration.(map[string]interface{}), nil, dockerHost, outputStream)
			if err == nil {
				dockerEndpoint := ""
				if dockerHost != nil {
					dockerEndpoint = dockerHost.Endpoint()
				}
				clusters = &contexts.Clusters{DockerHostEndpoint: dockerEndpoint}
				clusters.Add(router)
			}
		} else {
			errorMessage := "Invalid cluster configuration - \"cluster\" config missing"
			err = errors.New(errorMessage)
		}

		return clusters, err
	}
}

func parseCluster(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}, *contexts.Clusters, *docker_client.DockerHost, io.Writer) (*contexts.Cluster, error) {
	return func(clusterConfiguration map[string]interface{}, clusters *contexts.Clusters, dockerHost *docker_client.DockerHost, outputStream io.Writer) (*contexts.Cluster, error) {
		var (
			err                            error
			backendAddresses               []*contexts.BackendAddress
			dockerConfigurations           []*docker_client.DockerConfig
			version                        string
			sessionTimeout                 int64
			percentageTransitionPerRequest float64
			mode                           contexts.TransitionMode
			uuidValue                      uuid.UUID
			highestVersionCluster           *contexts.Cluster
		)

		if clusters != nil {
			highestVersionCluster = clusters.GetByVersionOrder(0)
		}

		uuidConfig := clusterConfiguration["uuid"]
		if uuidConfig != nil {
			uuidValue = uuid.Parse(uuidConfig.(string))
		} else {
			uuidValue = uuidGenerator()
		}

		versionConfig := clusterConfiguration["version"]
		if versionConfig != nil {
			floatVersion, isFloat := versionConfig.(float64)
			if isFloat {
				version = fmt.Sprintf("%.2f", floatVersion)
			} else {
				version = fmt.Sprintf("%s", versionConfig)
			}
		} else {
			version = "0.0"
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

		if err == nil {
			serversConfiguration := clusterConfiguration["servers"]
			containersConfiguration := clusterConfiguration["containers"]
			if serversConfiguration != nil {
				backendAddresses, err = parseServers(serversConfiguration)
			} else if containersConfiguration != nil {
				dockerConfigurations, backendAddresses, err = parseContainers(containersConfiguration, dockerHost)
			} else {
				errorMessage := "Invalid cluster configuration - \"cluster\" must contain \"servers\" or \"containers\" list"
				err = errors.New(errorMessage)
			}
		}

		if err == nil {
			if highestVersionCluster != nil && mode == contexts.ConcurrentMode {
				for _, existingBackendAddress := range highestVersionCluster.BackendAddresses {
					for _, newBackendAddress := range backendAddresses {
						if existingBackendAddress.Equals(newBackendAddress) {
							errorMessage := fmt.Sprintf("Invalid cluster configuration - new cluster has a conflicting address [%s] with existing highest version cluster [%s] (this is not allowed in \"CONCURRENT\" mode)", newBackendAddress, newBackendAddress)
							err = errors.New(errorMessage)
						}
					}
				}
			}
		}

		if err == nil {
			for _, dockerConfiguration := range dockerConfigurations {
				var dockerClient *docker_client.DockerClient
				dockerHost := dockerHost.Endpoint()
				if dockerConfiguration.DockerHost != nil && len(dockerConfiguration.DockerHost.Endpoint()) > 0 {
					dockerHost = dockerConfiguration.DockerHost.Endpoint()
				}
				dockerClient, err = docker_client.NewDockerClient(dockerHost)
				if err == nil {
					_, err = dockerClient.CreateServerFromContainer(dockerConfiguration, outputStream)
				}
			}
		}

		if err == nil {
			return &contexts.Cluster{BackendAddresses: backendAddresses, DockerConfigurations: dockerConfigurations, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: sessionTimeout, PercentageTransitionPerRequest: percentageTransitionPerRequest, Mode: mode, Version: version}, err
		} else {
			return nil, err
		}
	}
}

func parseServers(serversConfiguration interface{}) ([]*contexts.BackendAddress, error) {
	var (
		err                            error
		connection                     *net.TCPAddr
		backendAddresses               []*contexts.BackendAddress
	)

	servers, converted := serversConfiguration.([]interface{})
	if converted && len(servers) > 0 {
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
	} else {
		errorMessage := "Invalid cluster configuration - \"servers\" list must contain at least one entry"
		err = errors.New(errorMessage)
	}

	return backendAddresses, err
}


func parseContainers(containersConfiguration interface{}, dockerHost *docker_client.DockerHost) ([]*docker_client.DockerConfig, []*contexts.BackendAddress, error) {

	if dockerHost == nil {
		errorMessage := "Invalid docker host configuration - \"dockerHost\" must be provided when \"containers\" are specified"
		return nil, nil, errors.New(errorMessage)
	}

	var (
		err                   error
		connection            *net.TCPAddr
		backendAddresses      []*contexts.BackendAddress
		backendAddressesIndex int = 0
		dockerConfigurations  []*docker_client.DockerConfig
		marshaledBytes          []byte
	)

	containers, converted := containersConfiguration.([]interface{})
	if converted && len(containers) > 0 {
		dockerConfigurations = make([]*docker_client.DockerConfig, len(containers))
		backendAddresses = make([]*contexts.BackendAddress, len(containers))
		for index, containerConfig := range containers {
			container := containerConfig.(map[string]interface{})
			var dockerConfig = &docker_client.DockerConfig{}
			if marshaledBytes, err = json.Marshal(container); err == nil {
				if err = json.Unmarshal(marshaledBytes, dockerConfig); err == nil {
					if len(dockerConfig.Image) > 0 {
						dockerConfigurations[index] = dockerConfig
						if container["portToProxy"] != nil {
							portToProxy, isString := container["portToProxy"].(string)
							if !isString {
								portToProxy = fmt.Sprintf("%v", container["portToProxy"])
							}
							if len(portToProxy) > 0 {
								if dockerConfig.HasPortExposed(portToProxy) {
									dockerHostIp := dockerHost.Ip
									if dockerConfig.DockerHost != nil && len(dockerConfig.DockerHost.Endpoint()) > 0 {
										dockerHostIp = dockerConfig.DockerHost.Ip
									}

									connection, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", dockerHostIp, portToProxy))
									if err != nil {
										errorMessage := "Invalid container address [" + fmt.Sprintf("%s:%s", dockerHostIp, portToProxy) + "] - " + err.Error()
										err = errors.New(errorMessage)
									} else {
										backendAddresses[backendAddressesIndex] = &contexts.BackendAddress{Address: connection, Host: dockerHostIp, Port: portToProxy}
										backendAddressesIndex++
									}
								} else {
									errorMessage := "Invalid container configuration - port specified in \"portToProxy\" must be exposed by container in \"portBindings\""
									err = errors.New(errorMessage)
								}
							}
						}
					} else {
						errorMessage := "Invalid container configuration - no \"image\" specified"
						err = errors.New(errorMessage)
					}
				}
			}
		}
	} else {
		errorMessage := "Invalid cluster configuration - \"containers\" list must contain at least one entry"
		err = errors.New(errorMessage)
	}

	return dockerConfigurations, backendAddresses[0:backendAddressesIndex], err
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
		clusterMap := map[string]interface{}{"uuid": cluster.Uuid.String(), "servers": serversConfig, "version": cluster.Version, "upgradeTransition": upgradeTransition}
		if cluster.DockerConfigurations != nil {
			clusterMap["containers"] = cluster.DockerConfigurations
		}
		jsonConfig = map[string]interface{}{"cluster": clusterMap}
	}

	return jsonConfig
}

// ==== PARSE CONFIG - END
