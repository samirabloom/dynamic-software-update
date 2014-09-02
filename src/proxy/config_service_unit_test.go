package proxy

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
	"proxy/contexts"
	"proxy/docker_client"
)

func Test_Config_PUT_With_Valid_Json_Object(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":1024}, {\"hostname\":\"127.0.0.1\", \"port\":1025}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.1\"}}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}, URL: &url.URL{}}
		actualRouteContexts   = &contexts.Clusters{}
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRouteContexts = &contexts.Clusters{}
	)
	expectedRouteContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 1, Mode: contexts.SessionMode, Version: "1.1"})

	// when
	PUTHandler(uuidGenerator, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_With_Valid_Json_Object_In_Version_Order(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte1             = []byte("{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":1011}], \"upgradeTransition\":{\"sessionTimeout\":3}, \"version\": \"1.1\"}}")
		bodyByte2             = []byte("{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":1009}], \"upgradeTransition\":{\"sessionTimeout\":2}, \"version\": \"0.9\"}}")
		bodyByte3             = []byte("{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":1015}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.5\"}}")
		uuid1                 = uuid.Parse("1127596f-1034-11e4-8334-600308a82411")
		uuid2                 = uuid.Parse("0927596f-1034-11e4-8334-600308a82409")
		uuid3                 = uuid.Parse("1527596f-1034-11e4-8334-600308a82415")
		server1, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1011")
		server2, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1009")
		server3, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1015")
		expectedRouteContexts = &contexts.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*contexts.Cluster)}
		cluster1              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server1, Host: "127.0.0.1", Port: "1011"}}, RequestCounter: -1, Uuid: uuid1, SessionTimeout: 3, Mode: contexts.SessionMode, Version: "1.1"}
		cluster2              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server2, Host: "127.0.0.1", Port: "1009"}}, RequestCounter: -1, Uuid: uuid2, SessionTimeout: 2, Mode: contexts.SessionMode, Version: "0.9"}
		cluster3              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server3, Host: "127.0.0.1", Port: "1015"}}, RequestCounter: -1, Uuid: uuid3, SessionTimeout: 1, Mode: contexts.SessionMode, Version: "1.5"}
		actualRouteContexts   = &contexts.Clusters{}
	)

	// added in order
	expectedRouteContexts.ContextsByVersion.PushFront(cluster2) // 0.9
	expectedRouteContexts.ContextsByVersion.PushFront(cluster1) // 1.1 -> 0.9
	expectedRouteContexts.ContextsByVersion.PushFront(cluster3) // 1.5 -> 1.1 -> 0.9

	// added with key
	expectedRouteContexts.ContextsByID[uuid1.String()] = cluster1
	expectedRouteContexts.ContextsByID[uuid2.String()] = cluster2
	expectedRouteContexts.ContextsByID[uuid3.String()] = cluster3

	// when
	PUTHandler(func() uuid.UUID { return uuid1 }, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte1}, URL: &url.URL{}})
	PUTHandler(func() uuid.UUID { return uuid2 }, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte2}, URL: &url.URL{}})
	PUTHandler(func() uuid.UUID { return uuid3 }, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte3}, URL: &url.URL{}})

	// then
	assertion.AssertDeepEqual("Correct Object Added To Clusters In Version Order", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_With_Valid_Cluster_Configuration_Object(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{\"cluster\": {\"servers\": []}}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}, URL: &url.URL{}}
		actualRouteContexts   = &contexts.Clusters{}
		expectedRouteContexts = &contexts.Clusters{}
		expectedResponseBody  = []byte("Error parsing cluster configuration - Invalid cluster configuration - \"servers\" list must contain at least one entry\n")
	)

	// when
	PUTHandler(uuidGenerator, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusBadRequest, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, expectedResponseBody, responseWriter.WrittenBodyBytes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_When_Invalid_JSON(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{invalid}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}, URL: &url.URL{}}
		actualRouteContexts   = &contexts.Clusters{}
		expectedRouteContexts = &contexts.Clusters{}
		expectedResponseBody  = []byte("Error invalid character 'i' looking for beginning of object key string while decoding json {invalid}\n")
	)

	// when
	PUTHandler(uuidGenerator, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusBadRequest, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, expectedResponseBody, responseWriter.WrittenBodyBytes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_When_Empty_JSON(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}, URL: &url.URL{}}
		actualRouteContexts   = &contexts.Clusters{}
		expectedRouteContexts = &contexts.Clusters{}
		expectedResponseBody  = []byte("Invalid cluster configuration - \"cluster\" config missing\n")
	)

	// when
	PUTHandler(uuidGenerator, &docker_client.DockerHost{})(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusBadRequest, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, expectedResponseBody, responseWriter.WrittenBodyBytes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_GET_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		uuidValue            = uuidGenerator()
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts        = &contexts.Clusters{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("{\n    \"cluster\": {\n        \"servers\": [\n            {\n                \"hostname\": \"127.0.0.1\",\n                \"port\": 1024\n            },\n            {\n                \"hostname\": \"127.0.0.1\",\n                \"port\": 1025\n            }\n        ],\n        \"upgradeTransition\": {\n            \"mode\": \"SESSION\",\n            \"sessionTimeout\": 1\n        },\n        \"uuid\": \"" + uuidValue.String() + "\",\n        \"version\": \"1.1\"\n    }\n}")
		request              = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
	)
	routeContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: 1, Mode: contexts.SessionMode, Version: "1.1"})

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusOK, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, string(expectedResponseBody), string(responseWriter.WrittenBodyBytes[0]))
}

func Test_Config_GET_With_Non_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts        = &contexts.Clusters{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("404 page not found\n")
		request              = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
	)
	routeContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne}, &contexts.BackendAddress{Address: serverTwo}}, RequestCounter: -1, Uuid: uuidGenerator(), Version: "1.1"})

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, string(expectedResponseBody), string(responseWriter.WrittenBodyBytes[0]))
}

func Test_Config_GET_With_No_UUID(testCtx *testing.T) {
	// given
	var (
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		responseWriter       = mock.NewMockResponseWriter()
		uuid1                = uuid.Parse("1127596f-1034-11e4-8334-600308a82411")
		uuid2                = uuid.Parse("0927596f-1034-11e4-8334-600308a82409")
		uuid3                = uuid.Parse("1527596f-1034-11e4-8334-600308a82415")
		server1, _           = net.ResolveTCPAddr("tcp", "127.0.0.1:1011")
		server2, _           = net.ResolveTCPAddr("tcp", "127.0.0.1:1009")
		server3, _           = net.ResolveTCPAddr("tcp", "127.0.0.1:1015")
		cluster1             = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server1, Host: "127.0.0.1", Port: "1011"}}, RequestCounter: -1, Uuid: uuid1, SessionTimeout: 1, Mode: contexts.SessionMode, Version: "1.1"}
		cluster2             = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server2, Host: "127.0.0.1", Port: "1009"}}, RequestCounter: -1, Uuid: uuid2, SessionTimeout: 2, Mode: contexts.SessionMode, Version: "0.9"}
		cluster3             = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server3, Host: "127.0.0.1", Port: "1015"}}, RequestCounter: -1, Uuid: uuid3, SessionTimeout: 3, Mode: contexts.InstantMode, Version: "1.5"}
		routeContexts        = &contexts.Clusters{}
		expectedResponseBody = []byte("[\n    " +
			"{\n        \"cluster\": {\n            \"servers\": [\n                {\n                    \"hostname\": \"127.0.0.1\",\n                    \"port\": 1015\n                }\n            ],\n            \"upgradeTransition\": {\n                \"mode\": \"INSTANT\"\n            },\n            \"uuid\": \"" + uuid3.String() + "\",\n            \"version\": \"1.5\"\n        }\n    },\n    " +
			"{\n        \"cluster\": {\n            \"servers\": [\n                {\n                    \"hostname\": \"127.0.0.1\",\n                    \"port\": 1011\n                }\n            ],\n            \"upgradeTransition\": {\n                \"mode\": \"SESSION\",\n                \"sessionTimeout\": 1\n            },\n            \"uuid\": \"" + uuid1.String() + "\",\n            \"version\": \"1.1\"\n        }\n    },\n    " +
			"{\n        \"cluster\": {\n            \"servers\": [\n                {\n                    \"hostname\": \"127.0.0.1\",\n                    \"port\": 1009\n                }\n            ],\n            \"upgradeTransition\": {\n                \"mode\": \"SESSION\",\n                \"sessionTimeout\": 2\n            },\n            \"uuid\": \"" + uuid2.String() + "\",\n            \"version\": \"0.9\"\n        }\n    }\n" +
			"]")
		request              = &http.Request{URL: &url.URL{Path: "/server/"}}
	)
	routeContexts.Add(cluster1)
	routeContexts.Add(cluster2)
	routeContexts.Add(cluster3)

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusOK, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Response Body", testCtx, string(expectedResponseBody), string(responseWriter.WrittenBodyBytes[0]))
}

func Test_Config_DELETE_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts         = &contexts.Clusters{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
		expectedRouteContexts = &contexts.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*contexts.Cluster)}
	)
	routeContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne}, &contexts.BackendAddress{Address: serverTwo}}, RequestCounter: -1, Uuid: uuidGenerator(), Version: "1.1"})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Removed From Clusters", testCtx, expectedRouteContexts, routeContexts)
}


func Test_Config_DELETE_With_Existing_Object_Maintains_Order(testCtx *testing.T) {
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
		expectedRouteContexts = &contexts.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*contexts.Cluster)}
		cluster1              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server1}}, RequestCounter: -1, Uuid: uuid1, Version: "1.1"}
		cluster2              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server2}}, RequestCounter: -1, Uuid: uuid2, Version: "0.9"}
		cluster3              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: server3}}, RequestCounter: -1, Uuid: uuid3, Version: "1.5"}
		actualRouteContexts   = &contexts.Clusters{}
	)
	actualRouteContexts.Add(cluster1)
	actualRouteContexts.Add(cluster2)
	actualRouteContexts.Add(cluster3)

	// added in order
	expectedRouteContexts.ContextsByVersion.PushFront(cluster2) // 0.9
	expectedRouteContexts.ContextsByVersion.PushFront(cluster1) // 1.1 -> 0.9

	// added with key
	expectedRouteContexts.ContextsByID[uuid1.String()] = cluster1
	expectedRouteContexts.ContextsByID[uuid2.String()] = cluster2

	// when
	DeleteHandler(urlRegex)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Remove From Clusters In Version Order", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_DELETE_With_Non_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts         = &contexts.Clusters{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
		expectedRouteContexts = &contexts.Clusters{}
	)
	routeContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne}, &contexts.BackendAddress{Address: serverTwo}}, RequestCounter: -1, Uuid: uuidGenerator(), Version: "1.1"})
	expectedRouteContexts.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne}, &contexts.BackendAddress{Address: serverTwo}}, RequestCounter: -1, Uuid: uuidGenerator(), Version: "1.1"})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Removed From Clusters", testCtx, expectedRouteContexts, routeContexts)
}
