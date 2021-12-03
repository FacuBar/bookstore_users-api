package rest

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	t.Run("NoAuthorizationHeader", func(t *testing.T) {
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, nil))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "no authorization header was provided", err.Message())
	})

	t.Run("InvalidAuthorizationHeader", func(t *testing.T) {
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, nil))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "abc")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid authorization header format", err.Message())
	})

	t.Run("AuthorizationTypeNotSupported", func(t *testing.T) {
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, nil))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "notabearer abcd1234")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "authorization type not supported", err.Message())
	})

	t.Run("RestClientError", func(t *testing.T) {
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, &http.Client{}))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "Bearer token1234")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "couldn't verify session's validity", err.Message())
	})

	t.Run("InvalidAtNotAuthorized", func(t *testing.T) {
		testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(`{"message": "access_token not found","status": 404,"error": "not_found"}`))
		}))
		testServer.Listener.Close()
		l, _ := net.Listen("tcp", "127.0.0.1:8081")
		testServer.Listener = l
		testServer.Start()
		defer testServer.Close()

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, testServer.Client()))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "Bearer token1234")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "you are not logged in", err.Message())
	})

	t.Run("InvalidAtInternalServerError", func(t *testing.T) {
		testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "db error","status": 500,"error": "internal_server_error"}`))
		}))
		testServer.Listener.Close()
		l, _ := net.Listen("tcp", "127.0.0.1:8081")
		testServer.Listener = l
		testServer.Start()
		defer testServer.Close()

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) { c.Status(200) }, testServer.Client()))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "Bearer token1234")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.NotNil(t, err)
		assert.EqualValues(t, "couldn't verify session's validity", err.Message())
	})

	t.Run("NoError", func(t *testing.T) {
		// correct... mocks a succesfull call to the oauth-api
		testServer := correctTestServerResponse()
		defer testServer.Close()

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(resp)

		r.GET("/test", authenticate(func(c *gin.Context) {
			authorizedUser := c.MustGet(userPayloadKey).(userPayload)
			c.JSON(200, authorizedUser)
		}, testServer.Client()))

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Add("Authorization", "Bearer token1234")
		r.ServeHTTP(resp, c.Request)

		body, _ := ioutil.ReadAll(resp.Body)
		var uPayload userPayload
		json.Unmarshal(body, &uPayload)

		assert.EqualValues(t, 1, uPayload.Id)

		assert.EqualValues(t, 200, resp.Code)
	})
}
