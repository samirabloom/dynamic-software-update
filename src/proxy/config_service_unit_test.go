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
	"proxy/stages"
)

func Test_Config_PUT_With_Valid_Json_Object(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1024}, {\"ip\":\"127.0.0.1\", \"port\":1025}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		actualRouteContexts   = &stages.Clusters{}
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRouteContexts = &stages.Clusters{}
	)
	expectedRouteContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 1, Mode: stages.SessionMode, Version: 1.1})

	// when
	PUTHandler(uuidGenerator)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
}

func Test_Config_PUT_With_Valid_Json_Object_In_Version_Order(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte1             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1011}], \"upgradeTransition\":{\"sessionTimeout\":3}, \"version\": 1.1}}")
		bodyByte2             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1009}], \"upgradeTransition\":{\"sessionTimeout\":2}, \"version\": 0.9}}")
		bodyByte3             = []byte("{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1015}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.5}}")
		uuid1                 = uuid.Parse("1127596f-1034-11e4-8334-600308a82411")
		uuid2                 = uuid.Parse("0927596f-1034-11e4-8334-600308a82409")
		uuid3                 = uuid.Parse("1527596f-1034-11e4-8334-600308a82415")
		server1, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1011")
		server2, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1009")
		server3, _            = net.ResolveTCPAddr("tcp", "127.0.0.1:1015")
		expectedRouteContexts = &stages.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*stages.Cluster)}
		cluster1              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server1}, RequestCounter: -1, Uuid: uuid1, SessionTimeout: 3, Mode: stages.SessionMode, Version: 1.1}
		cluster2              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server2}, RequestCounter: -1, Uuid: uuid2, SessionTimeout: 2, Mode: stages.SessionMode, Version: 0.9}
		cluster3              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server3}, RequestCounter: -1, Uuid: uuid3, SessionTimeout: 1, Mode: stages.SessionMode, Version: 1.5}
		actualRouteContexts   = &stages.Clusters{}
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
	PUTHandler(func() uuid.UUID { return uuid1 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte1}})
	PUTHandler(func() uuid.UUID { return uuid2 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte2}})
	PUTHandler(func() uuid.UUID { return uuid3 })(actualRouteContexts, responseWriter, &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte3}})

	// then
	assertion.AssertDeepEqual("Correct Object Added To Clusters In Version Order", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_With_Valid_Cluster_Configuration_Object(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{\"cluster\": {\"servers\": []}}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		actualRouteContexts   = &stages.Clusters{}
		expectedRouteContexts = &stages.Clusters{}
	)

	// when
	PUTHandler(uuidGenerator)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusBadRequest, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_When_Invalid_JSON(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{invalid}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		actualRouteContexts   = &stages.Clusters{}
		expectedRouteContexts = &stages.Clusters{}
	)

	// when
	PUTHandler(uuidGenerator)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusBadRequest, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_PUT_When_Empty_JSON(testCtx *testing.T) {
	// given
	var (
		responseWriter        = mock.NewMockResponseWriter()
		bodyByte              = []byte("{}")
		request               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		actualRouteContexts   = &stages.Clusters{}
		expectedRouteContexts = &stages.Clusters{}
	)

	// when
	PUTHandler(uuidGenerator)(actualRouteContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Code", testCtx, 0, responseWriter.ResponseCodes[0])
	assertion.AssertDeepEqual("Correct Object Added To Clusters", testCtx, expectedRouteContexts, actualRouteContexts)
}

func Test_Config_GET_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		uuidValue            = uuidGenerator()
		urlRegex             = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts        = &stages.Clusters{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1024},{\"ip\":\"127.0.0.1\",\"port\":1025}],\"upgradeTransition\":{\"mode\":\"SESSION\",\"sessionTimeout\":1},\"uuid\":\"" + uuidValue.String() + "\",\"version\":1.1}}")
		request              = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
	)
	routeContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidValue, SessionTimeout: 1, Mode: stages.SessionMode, Version: 1.1})

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
		routeContexts        = &stages.Clusters{}
		responseWriter       = mock.NewMockResponseWriter()
		expectedResponseBody = []byte("404 page not found\n")
		request              = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
	)
	routeContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), Version: 1.1})

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Body", testCtx, string(expectedResponseBody), string(responseWriter.WritenBodyBytes[0]))
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
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
		cluster1             = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server1}, RequestCounter: -1, Uuid: uuid1, SessionTimeout: 1, Mode: stages.SessionMode, Version: 1.1}
		cluster2             = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server2}, RequestCounter: -1, Uuid: uuid2, SessionTimeout: 2, Mode: stages.SessionMode, Version: 0.9}
		cluster3             = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server3}, RequestCounter: -1, Uuid: uuid3, SessionTimeout: 3, Mode: stages.InstantMode, Version: 1.5}
		routeContexts        = &stages.Clusters{}
		expectedResponseBody = []byte("[" +
			"{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1015}],\"upgradeTransition\":{\"mode\":\"INSTANT\"},\"uuid\":\"" + uuid3.String() + "\",\"version\":1.5}}," +
			"{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1011}],\"upgradeTransition\":{\"mode\":\"SESSION\",\"sessionTimeout\":1},\"uuid\":\"" + uuid1.String() + "\",\"version\":1.1}}," +
			"{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1009}],\"upgradeTransition\":{\"mode\":\"SESSION\",\"sessionTimeout\":2},\"uuid\":\"" + uuid2.String() + "\",\"version\":0.9}}" +
			"]")
		request              = &http.Request{URL: &url.URL{Path: "/server/"}}
	)
	routeContexts.Add(cluster1)
	routeContexts.Add(cluster2)
	routeContexts.Add(cluster3)

	// when
	GETHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Response Body", testCtx, expectedResponseBody, responseWriter.WritenBodyBytes[0])
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusOK, responseWriter.ResponseCodes[0])
}

func Test_Config_DELETE_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		urlRegex              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		serverOne, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _          = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		routeContexts         = &stages.Clusters{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/" + uuidGenerator().String()}}
		expectedRouteContexts = &stages.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*stages.Cluster)}
	)
	routeContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), Version: 1.1})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Removed From Clusters", testCtx, expectedRouteContexts, routeContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusAccepted, responseWriter.ResponseCodes[0])
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
		expectedRouteContexts = &stages.Clusters{ContextsByVersion: list.New(), ContextsByID: make(map[string]*stages.Cluster)}
		cluster1              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server1}, RequestCounter: -1, Uuid: uuid1, Version: 1.1}
		cluster2              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server2}, RequestCounter: -1, Uuid: uuid2, Version: 0.9}
		cluster3              = &stages.Cluster{BackendAddresses: []*net.TCPAddr{server3}, RequestCounter: -1, Uuid: uuid3, Version: 1.5}
		actualRouteContexts   = &stages.Clusters{}
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
		routeContexts         = &stages.Clusters{}
		responseWriter        = mock.NewMockResponseWriter()
		request               = &http.Request{URL: &url.URL{Path: "/server/incorrect_uuid"}}
		expectedRouteContexts = &stages.Clusters{}
	)
	routeContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), Version: 1.1})
	expectedRouteContexts.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), Version: 1.1})

	// when
	DeleteHandler(urlRegex)(routeContexts, responseWriter, request)

	// then
	assertion.AssertDeepEqual("Correct Object Removed From Clusters", testCtx, expectedRouteContexts, routeContexts)
	assertion.AssertDeepEqual("Correct Response Code", testCtx, http.StatusNotFound, responseWriter.ResponseCodes[0])
}
