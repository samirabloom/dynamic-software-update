package proxy

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"proxy/log"
	"proxy/contexts"
	"os"
	"io"
	"net/url"
	"proxy/docker_client"
)

func ConfigServer(port int, routeContexts *contexts.Clusters, dockerHost *docker_client.DockerHost) {
	urlRegex := regexp.MustCompile("/configuration/cluster/([a-z0-9-]*){1}")
	http.ListenAndServe(":"+strconv.Itoa(port), &RegexpHandler{
			requestMappings: []*requestMapping{
				&requestMapping{pattern: regexp.MustCompile("/configuration/cluster"), method: "PUT", handler: PUTHandler(func() uuid.UUID {
					return uuid.NewUUID()
				}, dockerHost)},
				&requestMapping{pattern: urlRegex, method: "GET", handler: GETHandler(urlRegex)},
				&requestMapping{pattern: urlRegex, method: "DELETE", handler: DeleteHandler(urlRegex)}},
			routeContexts: routeContexts,
		})
}

type requestMapping struct {
	pattern *regexp.Regexp
	method string
	handler func(*contexts.Clusters, http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	requestMappings []*requestMapping
	routeContexts *contexts.Clusters
}

func (handler *RegexpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, requestMapping := range handler.requestMappings {
		if requestMapping.pattern.MatchString(request.URL.Path) && requestMapping.method == request.Method {
			requestMapping.handler(handler.routeContexts, writer, request)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(writer, request)
}

func PUTHandler(uuidGenerator func() uuid.UUID, dockerHost *docker_client.DockerHost) func(*contexts.Clusters, http.ResponseWriter, *http.Request) {
	return func(routeContexts *contexts.Clusters, writer http.ResponseWriter, request *http.Request) {
		body := make([]byte, 4096)
		size, _ := request.Body.Read(body)

		var jsonConfig map[string]interface{}
		err := json.Unmarshal(body[0:size], &jsonConfig)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error %s while decoding json %s", err.Error(), body[0:size]), http.StatusBadRequest)
		} else {
			clusterConfiguration := jsonConfig["cluster"]
			if clusterConfiguration != nil {
				var outputStream io.Writer = writer

				if request.URL != nil {
					parsedQueryString, queryParseErr := url.ParseQuery(request.URL.RawQuery)
					if queryParseErr == nil {
						logQueryParameter := parsedQueryString.Get("log")
						if logQueryParameter == "false" {
							outputStream = os.Stdout
						}
					}
				}
				cluster, err := parseCluster(uuidGenerator, false)(clusterConfiguration.(map[string]interface{}), routeContexts, dockerHost, outputStream)
				if err != nil {
					http.Error(writer, fmt.Sprintf("Error parsing cluster configuration - %s", err.Error()), http.StatusBadRequest)
				} else {
					routeContexts.Add(cluster)
					prettyBody, marshalErr := json.MarshalIndent(jsonConfig, "", "   ")
					if marshalErr == nil {
						log.LoggerFactory().Info(fmt.Sprintf("Received new cluster configuration:\n%s", prettyBody))
					} else {
						log.LoggerFactory().Info(fmt.Sprintf("Received new cluster configuration:\n%s", body[0:size]))
					}
					writer.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(writer, "%s", cluster.Uuid)
				}
			} else {
				http.Error(writer, "Invalid cluster configuration - \"cluster\" config missing", http.StatusBadRequest)
			}
		}
	}
}

func GETHandler(urlRegex *regexp.Regexp) func(*contexts.Clusters, http.ResponseWriter, *http.Request) {
	return func(routeContexts *contexts.Clusters, writer http.ResponseWriter, request *http.Request) {

		serverId := urlRegex.FindSubmatch([]byte(request.URL.Path))

		var (
			jsonBody []byte
			err error
		)

		uuidValue := string(serverId[1])
		if len(uuidValue) > 0 {
			routeContext := routeContexts.Get(uuid.Parse(uuidValue))
			if routeContext != nil {
				jsonBody, err = json.MarshalIndent(serialiseCluster(routeContext), "", "    ");
			} else {
				http.NotFound(writer, request)
				return
			}
		} else {
			index := 0
			var routeContextsJSON []map[string]interface{} = make([]map[string]interface{}, routeContexts.ContextsByVersion.Len())
			for element := routeContexts.ContextsByVersion.Front(); element != nil; element = element.Next() {
				routeContextsJSON[index] = serialiseCluster(element.Value.(*contexts.Cluster))
				index++
			}
			jsonBody, err = json.MarshalIndent(routeContextsJSON, "", "    ");
		}

		if err == nil {
			writer.WriteHeader(http.StatusOK)
			fmt.Fprintf(writer, "%s", jsonBody)
		}
	}
}


func DeleteHandler(urlRegex *regexp.Regexp) func(*contexts.Clusters, http.ResponseWriter, *http.Request) {
	return func(routeContexts *contexts.Clusters, writer http.ResponseWriter, request *http.Request) {
		serverId := urlRegex.FindSubmatch([]byte(request.URL.Path))

		uuid := uuid.Parse(string(serverId[1]))

		if routeContexts.Get(uuid) == nil {
			http.NotFound(writer, request)
		} else {
			routeContexts.Delete(uuid, writer)
			writer.WriteHeader(http.StatusAccepted)
		}
	}
}



