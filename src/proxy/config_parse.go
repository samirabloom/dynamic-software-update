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
	docker "github.com/fsouza/go-dockerclient"
	"proxy/docker_client"
	"io"
)

// ==== PARSE CONFIG - START

func loadConfig(configFile string, outputStream io.Writer) (*Proxy, error) {
	jsonData, err := readConfigFile(configFile)
	if err == nil {
		return parseConfigFile(jsonData, parseProxy, parseConfigService, parseClusters(func() uuid.UUID { return uuid.NewUUID() }, true), outputStream)
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

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseConfigService func(map[string]interface{}) (int, error), parseClusters func(map[string]interface{}, io.Writer) (*contexts.Clusters, error), outputStream io.Writer) (proxy *Proxy, err error) {
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
				clusters, clusterParseErr := parseClusters(jsonConfig, outputStream)
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

func parseClusters(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}, io.Writer) (*contexts.Clusters, error) {
	return func(jsonConfig map[string]interface{}, outputStream io.Writer) (*contexts.Clusters, error) {
		var (
			err      error
			router   *contexts.Cluster
			clusters *contexts.Clusters
		)

		clusterConfiguration := jsonConfig["cluster"]
		if clusterConfiguration != nil {
			router, err = parseCluster(uuidGenerator, initialCluster)(clusterConfiguration.(map[string]interface{}), outputStream)
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

func parseCluster(uuidGenerator func() uuid.UUID, initialCluster bool) func(map[string]interface{}, io.Writer) (*contexts.Cluster, error) {
	return func(clusterConfiguration map[string]interface{}, outputStream io.Writer) (*contexts.Cluster, error) {
		var (
			err                            error
			backendAddresses               []*contexts.BackendAddress
			version                        string
			sessionTimeout                 int64
			percentageTransitionPerRequest float64
			mode                           contexts.TransitionMode
			uuidValue                      uuid.UUID
		)

		serversConfiguration := clusterConfiguration["servers"]
		containersConfiguration := clusterConfiguration["containers"]
		if serversConfiguration != nil {
			backendAddresses, err = parseServers(serversConfiguration)
		} else if containersConfiguration != nil {
			backendAddresses, err = parseContainers(containersConfiguration, outputStream)
		} else {
			errorMessage := "Invalid cluster configuration - \"cluster\" must contain \"servers\" or \"containers\" list"
			err = errors.New(errorMessage)
		}

		if err == nil {
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
			return &contexts.Cluster{BackendAddresses: backendAddresses, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: sessionTimeout, PercentageTransitionPerRequest: percentageTransitionPerRequest, Mode: mode, Version: version}, err
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

/*
{
*	"image":"",                         // Image for container
*   "portToProxy":
	"version":

	"workingDir":"",  					// Working directory inside the container
	"entrypoint":"",  					// Overwrite the default ENTRYPOINT of the image
	"env":null,       					// Set environment variables
	"cmd":[                             // Set command executed when the container runs
		 ""
	],

	"hostname":"",   					// Container host name
	"volumes":{       					// Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)
		 "/tmp": {}
	},
	"volumesFrom":[                 	// Mount volumes from the specified container(s)
		 "parent",
		 "other:ro"
	],
	"exposedPorts":{  					// Expose a port from the container without publishing it to your host
		 "22/tcp": {}
	},
	"publishAllPorts":false,            // Publish all exposed ports to the host interfaces
*	"portBindings":{ "22/tcp": [{ "HostPort": "11022" }] },
	"links":["redis3:redis"],       	// Add link to another container in the form of name:alias

	"user":"",       					// Username or UID
	"memory":0,      					// Memory limit (format: <number><optional unit>, where unit = b, k, m or g)
	"cpuShares":0                   	// CPU shares (relative weight)
	"lxcConf":{"lxc.utsname":"docker"}  // (lxc exec-driver only) Add custom lxc options --lxc-conf="lxc.cgroup.cpuset.cpus = 0,1"
	"privileged":false                  // Give extended privileges to this container
}
 */


func parseContainers(containersConfiguration interface{}, outputStream io.Writer) ([]*contexts.BackendAddress, error) {
	var (
		err                 error
		connection          *net.TCPAddr
		backendAddresses    []*contexts.BackendAddress
		image               string                               = ""
		tag                 string                               = ""
		portToProxy         int64                                = 0
		workingDir          string                               = ""
		entrypoint          string                               = ""
		env                 string                               = ""
		cmd                 []string                             = nil
		hostname            string                               = ""
		volumes             []string                             = nil
		volumesFrom         []string                             = nil
		exposedPorts        map[docker.Port]struct{}             = nil
		publishAllPorts     bool                                 = false
		portBindings        map[docker.Port][]docker.PortBinding = nil
		links               []string                             = nil
		user                string                               = ""
		memory              int64                                = 0
		cpuShares           int64                                = 0
		lxcConf             []docker.KeyValuePair                = nil
		privileged          bool                                 = false
	)

	containers, converted := containersConfiguration.([]interface{})
	if converted && len(containers) > 0 {
		backendAddresses = make([]*contexts.BackendAddress, len(containers))
		for index := range containers {
			var container map[string]interface{} = containers[index].(map[string]interface{})

			imageConfig := container["image"]
			if imageConfig != nil {
				image = imageConfig.(string)
			}

			tagConfig := container["tag"]
			if tagConfig != nil {
				tag = tagConfig.(string)
			}

			portToProxyConfig := container["portToProxy"]
			if portToProxyConfig != nil {
				portToProxy = int64(portToProxyConfig.(float64))
			}

			workingDirConfig := container["workingDir"]
			if workingDirConfig != nil {
				workingDir = workingDirConfig.(string)
			}

			entrypointConfig := container["entrypoint"]
			if entrypointConfig != nil {
				entrypoint = entrypointConfig.(string)
			}

			envConfig := container["env"]
			if envConfig != nil {
				env = envConfig.(string)
			}

			cmdConfig := container["cmd"]
			if cmdConfig != nil {
				cmd = cmdConfig.([]string)
			}

			hostnameConfig := container["hostname"]
			if hostnameConfig != nil {
				hostname = hostnameConfig.(string)
			}

			volumesConfig := container["volumes"]
			if volumesConfig != nil {
				volumes = volumesConfig.([]string)
			}

			volumesFromConfig := container["volumesFrom"]
			if volumesFromConfig != nil {
				volumesFrom = volumesFromConfig.([]string)
			}

			exposedPortsConfig := container["exposedPorts"]
			if exposedPortsConfig != nil {
				exposedPorts = exposedPortsConfig.(map[docker.Port]struct {})
			}

			publishAllPortsConfig := container["publishAllPorts"]
			if publishAllPortsConfig != nil {
				publishAllPorts = publishAllPortsConfig.(bool)
			}

			portBindingsConfig := container["portBindings"]
			if portBindingsConfig != nil {
				portBindings = portBindingsConfig.(map[docker.Port][]docker.PortBinding)
			}

			linksConfig := container["links"]
			if linksConfig != nil {
				links = linksConfig.([]string)
			}

			userConfig := container["user"]
			if userConfig != nil {
				user = userConfig.(string)
			}

			memoryConfig := container["memory"]
			if memoryConfig != nil {
				memory = int64(memoryConfig.(float64))
			}

			cpuSharesConfig := container["cpuShares"]
			if cpuSharesConfig != nil {
				cpuShares = int64(cpuSharesConfig.(float64))
			}

			lxcConfConfig := container["lxcConf"]
			if lxcConfConfig != nil {
				lxcConf = lxcConfConfig.([]docker.KeyValuePair)
			}

			privilegedConfig := container["privileged"]
			if privilegedConfig != nil {
				privileged = privilegedConfig.(bool)
			}

			dockerConfig := &docker_client.DockerConfig{
				Image: image,
				Tag: tag,
				PortToProxy: portToProxy,
				WorkingDir: workingDir,
				Entrypoint: entrypoint,
				Env: env,
				Cmd: cmd,
				Hostname: hostname,
				Volumes: volumes,
				VolumesFrom: volumesFrom,
				ExposedPorts: exposedPorts,
				PublishAllPorts: publishAllPorts,
				PortBindings: portBindings,
				Links: links,
				User: user,
				Memory: memory,
				CpuShares: cpuShares,
				LxcConf: lxcConf,
				Privileged: privileged,
			}

			fmt.Printf("dockerConfig: %#v\n", dockerConfig)

			dockerClient, err := docker_client.NewDockerClient("http://192.168.50.5:2375")
			if err != nil {
				dockerClient.CreateServerFromContainer(dockerConfig, outputStream)
			}

			fmt.Printf("dockerClient: %#v\n", dockerClient)

			// todo
			// 1. add configuration for docker REST API host and port
			// 1. query ip address of container
			// 2. add portToProxy & ip address of container as backendAddresses

			connection, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", container["hostname"], container["port"]))
			if err != nil {
				errorMessage := "Invalid container address [" + fmt.Sprintf("%s:%v", container["hostname"], container["port"]) + "] - " + err.Error()
				err = errors.New(errorMessage)
			} else {
				backendAddresses[index] = &contexts.BackendAddress{Address: connection, Host: fmt.Sprintf("%s", container["hostname"]), Port: fmt.Sprintf("%v", container["port"])}
			}
		}
	} else {
		errorMessage := "Invalid cluster configuration - \"containers\" list must contain at least one entry"
		err = errors.New(errorMessage)
	}

	return backendAddresses, err
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
