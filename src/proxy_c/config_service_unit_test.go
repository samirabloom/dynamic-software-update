package proxy_c

import (
	"testing"
	"net/http"
	"regexp"
	"net/url"
	mock "util/test/mock"
	assertion "util/test/assertion"
	"net"
	"code.google.com/p/go-uuid/uuid"
	"container/list"
)

func Test_Config_PUT_With_Valid_Json_Object(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1024}, {\"ip\":\"127.0.0.1\", \"port\":1025}], \"version\": 1.1}}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		actualRouteContexts   = &RoutingContexts{}
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRouteContexts = &RoutingContexts{}
	)
	expectedRouteContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.1})

	// when
	PUTHandler(uuidGenerator)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Added To RoutingContexts", testCtx, expectedRouteContexts, actualRouteContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
}

func Test_Config_PUT_With_Valid_Json_Object_In_Version_Order(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte1             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1011}], \"version\": 1.1}}")
		bodyByte2             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1009}], \"version\": 0.9}}")
		bodyByte3             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1015}], \"version\": 1.5}}")
		uuid1                 = uuid.Parse("1127596f-1034-11e4-8334-600308a82411")
		uuid2                 = uuid.Parse("0927596f-1034-11e4-8334-600308a82409")
		uuid3                 = uuid.Parse("1527596f-1034-11e4-8334-600308a82415")
		server1, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1011")
		server2, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1009")
		server3, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1015")
		expectedRouteContexts = &RoutingContexts{contextsByVersion: list.New(), contextsByID: make(map[string]*RoutingContext)}
		routingContext1       = &RoutingContext{backendAddresses: []*net.TCPAddr{server1}, requestCounter: -1, uuid: uuid1, version: 1.1}
		routingContext2       = &RoutingContext{backendAddresses: []*net.TCPAddr{server2}, requestCounter: -1, uuid: uuid2, version: 0.9}
		routingContext3       = &RoutingContext{backendAddresses: []*net.TCPAddr{server3}, requestCounter: -1, uuid: uuid3, version: 1.5}
		actualRouteContexts   = &RoutingContexts{}
	)

	// added in order
	expectedRouteContexts.contextsByVersion.PushFront(routingContext2) // 0.9
	expectedRouteContexts.contextsByVersion.PushFront(routingContext1) // 1.1 -> 0.9
	expectedRouteContexts.contextsByVersion.PushFront(routingContext3) // 1.5 -> 1.1 -> 0.9

	// added with key
	expectedRouteContexts.contextsByID[uuid1.String()] = routingContext1
	expectedRouteContexts.contextsByID[uuid2.String()] = routingContext2
	expectedRouteContexts.contextsByID[uuid3.String()] = routingContext3

	// when
	PUTHandler(func() uuid.UUID { return uuid1 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte1}})
	PUTHandler(func() uuid.UUID { return uuid2 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte2}})
	PUTHandler(func() uuid.UUID { return uuid3 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte3}})

	// then
	assertion.AssertDeepEqual("Correct Object Added To RoutingContexts In Version Order", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_GET_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		uuidValue            = uuidGenerator()
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts        = &RoutingContexts{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1024},{\"ip\":\"127.0.0.1\",\"port\":1025}],\"uuid\":\"" + uuidValue.String() + "\",\"version\":1.1}}")
		request              = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
	)
	routeContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidValue, version: 1.1})

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Body", testCtx, expectedResponseBody, responseWriter.WritenBodyBytes[0])
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusOK, responseWriter.ResponseCodes[0])
}

func Test_Config_GET_With_Non_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts        = &RoutingContexts{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("404 page not found\n")
		request              = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
	)
	routeContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.1})

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Body", testCtx, string(expectedResponseBody), string(responseWriter.WritenBodyBytes[0]))
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
}

func Test_Config_DELETE_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts         = &RoutingContexts{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
		expectedRouteContexts = &RoutingContexts{contextsByVersion: list.New(), contextsByID: make(map[string]*RoutingContext)}
	)
	routeContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.1})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Removed From RoutingContexts", testCtx, expectedRouteContexts, routeContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
}


func Test_Config_DELETE_With_Existing_Object_Maintanes_Order(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		responseWriter        = mock.NewMockResponseWriter()
		uuid1                 = uuid.Parse("1127596f-1034-11e4-8334-600308a82411")
		uuid2                 = uuid.Parse("0927596f-1034-11e4-8334-600308a82409")
		uuid3                 = uuid.Parse("1527596f-1034-11e4-8334-600308a82415")
		request               = &http.Request{URL: &url.URL{Path: "/server/" + uuid3.String()}}
		server1, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1011")
		server2, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1009")
		server3, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1015")
		expectedRouteContexts = &RoutingContexts{contextsByVersion: list.New(), contextsByID: make(map[string]*RoutingContext)}
		routingContext1       = &RoutingContext{backendAddresses: []*net.TCPAddr{server1}, requestCounter: -1, uuid: uuid1, version: 1.1}
		routingContext2       = &RoutingContext{backendAddresses: []*net.TCPAddr{server2}, requestCounter: -1, uuid: uuid2, version: 0.9}
		routingContext3       = &RoutingContext{backendAddresses: []*net.TCPAddr{server3}, requestCounter: -1, uuid: uuid3, version: 1.5}
		actualRouteContexts   = &RoutingContexts{}
	)
	actualRouteContexts.Add(routingContext1)
	actualRouteContexts.Add(routingContext2)
	actualRouteContexts.Add(routingContext3)

	// added in order
	expectedRouteContexts.contextsByVersion.PushFront(routingContext2) // 0.9
	expectedRouteContexts.contextsByVersion.PushFront(routingContext1) // 1.1 -> 0.9

	// added with key
	expectedRouteContexts.contextsByID[uuid1.String()] = routingContext1
	expectedRouteContexts.contextsByID[uuid2.String()] = routingContext2

	// when
	DeleteHandler(urlRegex)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Remove From RoutingContexts In Version Order", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_DELETE_With_Non_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts         = &RoutingContexts{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
		expectedRouteContexts = &RoutingContexts{}
	)
	routeContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.1})
	expectedRouteContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.1})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Removed From RoutingContexts", testCtx, expectedRouteContexts, routeContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
}
