package proxy_config

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"time"
)

func main() {
	Server(8080)
	time.Sleep(100 * time.Millisecond)
}

func Server(port int) {
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
		// retriving json bodies
		body := make([]byte, 1024)
		size, _ := r.Body.Read(body)

		var jsonObject interface {}
		err := json.Unmarshal(body[0:size], &jsonObject)
		if err != nil {
			fmt.Println("Error decoding json request", err.Error())
		}
		id := uuidGenerator()
		jsonObject.(map[string]interface{})["id"] = id
		jsonObjectMaps[id] = jsonObject

		//	for key, value := range jsonObjectMaps {
		//		fmt.Printf("\njson with id \n %s is: \n %s\n\n", key, value)
		//	}

		// parse received json body
//		m := jsonObject.(map[string]interface{})
//		parseJsonBody(m)

		fmt.Fprintf(w, "%s", id)
	}
}

func parseJsonBody(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
		for i, u := range vv {
			fmt.Printf("[%s], [%s]", i, u)
		}
		case map[string]interface{}:
			fmt.Printf("%s is a map: ", k)
		for i, u := range vv {
			fmt.Printf("\n\t[%s], [%s]", i, u)
		}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
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
		//		println("I am in the get function")
		serverId := urlRegex.FindSubmatch([]byte(r.URL.Path))
		if len(serverId) >= 2 {
			//fmt.Printf("dynsoftup value is: %s\n", string(serverId[1]))
		}
		if jsonObjectMaps[string(serverId[1])] == nil {
			http.NotFound(w, r)
		} else {
			delete(jsonObjectMaps, string(serverId[1]))
			response := http.StatusText(202)
			fmt.Fprintf(w, "%s", response)
		}
	}
}



