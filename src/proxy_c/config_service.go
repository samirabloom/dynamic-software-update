package proxy_c

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"errors"
)

func ConfigServer(port float64, routeContexts *RoutingContexts) {
	urlRegex := regexp.MustCompile("/server/([a-z0-9-]*){1}")
	http.ListenAndServe(":"+strconv.Itoa(int(port)), &RegexpHandler{
		requestMappings: []*requestMapping{
		&requestMapping{pattern: regexp.MustCompile("/server"), method: "PUT", handler: PUTHandler(func() uuid.UUID {
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
	handler func(*RoutingContexts, http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	requestMappings []*requestMapping
	routeContexts *RoutingContexts
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

func PUTHandler(uuidGenerator func() uuid.UUID) func(*RoutingContexts, http.ResponseWriter, *http.Request) {
	return func(routeContexts *RoutingContexts, writer http.ResponseWriter, request *http.Request) {
		body := make([]byte, 1024)
		size, _ := request.Body.Read(body)

		var jsonConfig map[string]interface{}
		err := json.Unmarshal(body[0:size], &jsonConfig)
		if err != nil {
			fmt.Printf("Error decoding json request:\n\t%s\n", err.Error())
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			if jsonConfig == nil {
				fmt.Printf("Error parsing cluster configuration invalid JSON\n")
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			} else {
				clusterConfiguration := jsonConfig["cluster"]
				if clusterConfiguration != nil {
					routingContext, err := parseRoutingContext(uuidGenerator)(clusterConfiguration.(map[string]interface {}))
					if err != nil {
						fmt.Printf("Error parsing cluster configuration:\n\t%s\n", err.Error())
						http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					} else {
						routeContexts.Add(routingContext)
						loggerFactory().Info(fmt.Sprintf("Received new cluster configuration:\n%s", body[0:size]))
						writer.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(writer, "%s", routingContext.uuid)
					}
				} else {
					errorMessage := "Invalid cluster configuration - \"cluster\" config missing"
					loggerFactory().Error(errorMessage)
					err = errors.New(errorMessage)
				}
			}
		}
	}
}

func GETHandler(urlRegex *regexp.Regexp) func(*RoutingContexts, http.ResponseWriter, *http.Request) {
	return func(routeContexts *RoutingContexts, writer http.ResponseWriter, request *http.Request) {

		serverId := urlRegex.FindSubmatch([]byte(request.URL.Path))
		if len(serverId) >= 2 {
		}

		uuidValue := string(serverId[1])
		routeContext := routeContexts.Get(uuid.Parse(uuidValue))
		jsonBody, err := json.Marshal(serialiseRoutingContext(routeContext));

		if err != nil {
			fmt.Println("Error encoding json object %s", uuidValue)
		} else if routeContext == nil {
			http.NotFound(writer, request)
		} else {
			fmt.Fprintf(writer, "%s", jsonBody)
			writer.WriteHeader(http.StatusOK)
		}
	}
}


func DeleteHandler(urlRegex *regexp.Regexp) func(*RoutingContexts, http.ResponseWriter, *http.Request) {
	return func(routeContexts *RoutingContexts, writer http.ResponseWriter, request *http.Request) {
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



