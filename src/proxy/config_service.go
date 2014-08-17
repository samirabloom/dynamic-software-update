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
)

func ConfigServer(port int, routeContexts *contexts.Clusters) {
	urlRegex := regexp.MustCompile("/configuration/cluster/([a-z0-9-]*){1}")
	http.ListenAndServe(":"+strconv.Itoa(port), &RegexpHandler{
		requestMappings: []*requestMapping{
		&requestMapping{pattern: regexp.MustCompile("/configuration/cluster"), method: "PUT", handler: PUTHandler(func() uuid.UUID {
			return uuid.NewUUID()
		})},
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

func PUTHandler(uuidGenerator func() uuid.UUID) func(*contexts.Clusters, http.ResponseWriter, *http.Request) {
	return func(routeContexts *contexts.Clusters, writer http.ResponseWriter, request *http.Request) {
		body := make([]byte, 1024)
		size, _ := request.Body.Read(body)

		var jsonConfig map[string]interface{}
		err := json.Unmarshal(body[0:size], &jsonConfig)
		if err != nil {
			fmt.Fprintf(writer, "Error decoding json request - %s", err.Error())
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			clusterConfiguration := jsonConfig["cluster"]
			if clusterConfiguration != nil {
				cluster, err := parseCluster(uuidGenerator, false)(clusterConfiguration.(map[string]interface{}))
				if err != nil {
					fmt.Fprintf(writer, "Error parsing cluster configuration - %s", err.Error())
					http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				} else {
					routeContexts.Add(cluster)
					log.LoggerFactory().Info(fmt.Sprintf("Received new cluster configuration:\n%s", body[0:size]))
					writer.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(writer, "%s", cluster.Uuid)
				}
			} else {
				fmt.Fprintf(writer, "Invalid cluster configuration - \"cluster\" config missing")
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
				jsonBody, err = json.Marshal(serialiseCluster(routeContext));
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
			jsonBody, err = json.Marshal(routeContextsJSON);
		}

		if err == nil {
			fmt.Fprintf(writer, "%s", jsonBody)
			writer.WriteHeader(http.StatusOK)
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
			routeContexts.Delete(uuid)
			writer.WriteHeader(http.StatusAccepted)
		}
	}
}



