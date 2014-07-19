package proxy_config

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"time"
	"proxy_c"
)

func main() {
	ConfigServer(8080, nil)
	time.Sleep(100 * time.Millisecond)
}

func ConfigServer(port int, contexts *proxy_c.RoutingContexts) {
	fmt.Println("Starting server " + strconv.Itoa(port) + " ....")
	urlRegex := regexp.MustCompile("/server/([a-z0-9-]*){1}")
	http.ListenAndServe(":"+strconv.Itoa(port), &RegexpHandler{
		routes: []*route{
		&route{pattern: regexp.MustCompile("/server"), method: "PUT", handler: PUTHandler(func() string {
			return uuid.NewUUID().String()
		})},
		&route{pattern: urlRegex, method: "GET", handler: GETHandler(urlRegex)},
		&route{pattern: urlRegex, method: "DELETE", handler: DeleteHandler(urlRegex)},
	},
		jsonObjectMaps: make(map[string]interface{}),
	},
	)
}

type route struct {
	pattern *regexp.Regexp
	method string
	handler func(map[string]interface{}, http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	routes []*route
	jsonObjectMaps map[string]interface{}
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.method == r.Method {
			route.handler(h.jsonObjectMaps, w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func PUTHandler(uuidGenerator func() string) func(map[string]interface{}, http.ResponseWriter, *http.Request) {
	return func(jsonObjectMaps map[string]interface{}, w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 1024)
		size, _ := r.Body.Read(body)

		var jsonObject interface{}
		err := json.Unmarshal(body[0:size], &jsonObject)
		if err != nil {
			fmt.Println("Error decoding json request", err.Error())
		}
		id := uuidGenerator()
		jsonObject.(map[string]interface{})["id"] = id
		jsonObjectMaps[id] = jsonObject

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "%s", id)
	}
}

func GETHandler(urlRegex *regexp.Regexp) func(map[string]interface{}, http.ResponseWriter, *http.Request) {
	return func(jsonObjectMaps map[string]interface{}, writer http.ResponseWriter, request *http.Request) {
		//		println("I am in the get function")
		serverId := urlRegex.FindSubmatch([]byte(request.URL.Path))
		if len(serverId) >= 2 {
		}

		jsonBody, err := json.Marshal(jsonObjectMaps[string(serverId[1])]);

		if err != nil {
			fmt.Println("Error encoding json object\n%s :", jsonObjectMaps[string(serverId[1])])
		} else if jsonObjectMaps[string(serverId[1])] == nil {
			http.NotFound(writer, request)
		} else {
			fmt.Fprintf(writer, "%s", jsonBody)
		}
	}
}



func DeleteHandler(urlRegex *regexp.Regexp) func(map[string]interface{}, http.ResponseWriter, *http.Request) {
	return func(jsonObjectMaps map[string]interface{}, w http.ResponseWriter, r *http.Request) {
		serverId := urlRegex.FindSubmatch([]byte(r.URL.Path))

		if jsonObjectMaps[string(serverId[1])] == nil {
			http.NotFound(w, r)
		} else {
			delete(jsonObjectMaps, string(serverId[1]))
			w.WriteHeader(http.StatusAccepted)
		}
	}
}



