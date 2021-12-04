package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/stretchr/testify/assert"
)

// Mocking services
type usersServiceMock struct {
}

var (
	funcGetUser  func(int64) (*domain.User, rest_errors.RestErr)
	funcRegister func(*domain.User) rest_errors.RestErr
	funcUpdate   func(*domain.User, bool) rest_errors.RestErr
	funcLogin    func(string, string) (*domain.User, rest_errors.RestErr)
)

func (m *usersServiceMock) GetUser(id int64) (*domain.User, rest_errors.RestErr) {
	return funcGetUser(id)
}
func (m *usersServiceMock) Register(user *domain.User) rest_errors.RestErr {
	return funcRegister(user)
}
func (m *usersServiceMock) Update(user *domain.User, isAdmin bool) rest_errors.RestErr {
	return funcUpdate(user, isAdmin)
}
func (m *usersServiceMock) Login(email string, password string) (*domain.User, rest_errors.RestErr) {
	return funcLogin(email, password)
}

var usm ports.UsersService = &usersServiceMock{}

var userTest = domain.User{
	Id:          1,
	FirstName:   "Oscar",
	LastName:    "Isaac",
	Email:       "oscaac@gmail.com",
	Password:    "$2a$10$jRL.gYiodDnwcOBErnDfuu5044h40PM3ZOAOzit6O4RIL9wG24xJ6", //password
	DateCreated: "2006-01-02 15:04:05",
	Status:      "active",
	Role:        "user",
}

func TestRegisterUser(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcRegister = func(u *domain.User) rest_errors.RestErr {
			return nil
		}

		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"first_name":"Oscar","last_name":"Isaac","email": "oscaac@gmail.com","password":"somepass","confirm_password":"somepass"}`))
		server.srv.Handler.ServeHTTP(w, req)

		assert.EqualValues(t, http.StatusCreated, w.Code)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"first_name":1}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		resp, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, w.Code)
		assert.EqualValues(t, "invalid request", resp.Message())
	})

	t.Run("PasswordsNotEqual", func(t *testing.T) {
		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"password":"somepass","confirm_password":"notthesamepass"}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		resp, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, w.Code)
		assert.EqualValues(t, "passwords are not equal", resp.Message())
	})

	t.Run("ServiceError", func(t *testing.T) {
		funcRegister = func(u *domain.User) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("error while trying to register, try again later")
		}

		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"first_name":"Oscar","last_name":"Isaac","email": "oscaac@gmail.com","password":"somepass","confirm_password":"somepass"}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		resp, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusInternalServerError, resp.Status())
		assert.EqualValues(t, "error while trying to register, try again later", resp.Message())
	})
}

func TestGetUser(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcGetUser = func(i int64) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}

		// correct... mocks a succesfull call to the oauth-api made inside of the middleware
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		var user domain.User
		json.Unmarshal(body, &user)

		assert.EqualValues(t, http.StatusOK, w.Code)
		assert.EqualValues(t, userTest.Email, user.Email)
	})

	t.Run("InvalidId", func(t *testing.T) {

		// correct... mocks a succesfull call to the oauth-api made inside of the middleware
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/abc", nil)
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, err.Status())
		assert.EqualValues(t, "user id not valid", err.Message())
	})

	t.Run("LackPermissions", func(t *testing.T) {
		// correct... mocks a succesfull call to the oauth-api made inside of the middleware
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/2", nil)
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusUnauthorized, err.Status())
		assert.EqualValues(t, "you don't have the permissions to access this resource", err.Message())
	})

	t.Run("ServiceError", func(t *testing.T) {
		funcGetUser = func(i int64) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}

		// correct... mocks a succesfull call to the oauth-api made inside of the middleware
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "user not found", err.Message())
	})
}

func TestLogin(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcLogin = func(s1, s2 string) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}

		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(`{"email": "oscaac@gmail.com","password":"123456"}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		var user domain.User
		json.Unmarshal(body, &user)

		assert.EqualValues(t, http.StatusOK, w.Code)
		assert.EqualValues(t, userTest.Email, user.Email)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		funcLogin = func(s1, s2 string) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}

		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(`{"email": "oscaac@gmail.com","password":123456}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)
		fmt.Println(err.Message())

		assert.EqualValues(t, http.StatusBadRequest, err.Status())
		assert.EqualValues(t, "invalid request", err.Message())
	})

	t.Run("ServiceError", func(t *testing.T) {
		funcLogin = func(s1, s2 string) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewBadRequestError("invalid credentials")
		}

		server := NewServer(&http.Server{}, nil, nil, nil, nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(`{"email": "oscaac@gmail.com","password":"wrongpassword"}`))
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)
		fmt.Println(err.Message())

		assert.EqualValues(t, http.StatusBadRequest, err.Status())
		assert.EqualValues(t, "invalid credentials", err.Message())
	})
}

func TestUpdate(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcUpdate = func(u *domain.User, b bool) rest_errors.RestErr {
			u.FirstName = "Oscar"

			return nil
		}

		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"email": "oscar@newemail.com"}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		var user domain.User
		json.Unmarshal(body, &user)

		assert.EqualValues(t, http.StatusOK, w.Code)
		assert.EqualValues(t, "Oscar", user.FirstName)
		assert.EqualValues(t, "oscar@newemail.com", user.Email)
	})

	t.Run("InvalidId", func(t *testing.T) {
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/abc", strings.NewReader(`{"email": "oscar@newemail.com"}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		restErr, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, restErr.Status())
		assert.EqualValues(t, "user id not valid", restErr.Message())
	})

	t.Run("UnauthorizedSession", func(t *testing.T) {
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/2", strings.NewReader(`{"email": "oscar@newemail.com"}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		restErr, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusUnauthorized, restErr.Status())
		assert.EqualValues(t, "you don't have the permissions to perform this action", restErr.Message())
	})

	t.Run("InvalidReqBody", func(t *testing.T) {
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"email": 2}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		restErr, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, restErr.Status())
		assert.EqualValues(t, "invalid request body", restErr.Message())
	})

	t.Run("PasswordsNotEqual", func(t *testing.T) {
		testServer := correctTestServerResponse()
		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"password":"some","confirm_password":"any"}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		restErr, _ := rest_errors.NewRestErrorFromBytes(body)

		assert.EqualValues(t, http.StatusBadRequest, restErr.Status())
		assert.EqualValues(t, "passwords are not equal", restErr.Message())
	})

	t.Run("Error", func(t *testing.T) {
		testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{"access_token":"084a4a0f-92cc-46e6-9b57-1d2aed3c389e", "user_id":1, "user_role":"admin", "expires":1637510344}`))
		}))
		testServer.Listener.Close()
		l, _ := net.Listen("tcp", "127.0.0.1:8081")
		testServer.Listener = l
		testServer.Start()

		funcUpdate = func(u *domain.User, b bool) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("db error")
		}

		defer testServer.Close()

		server := NewServer(&http.Server{}, nil, nil, testServer.Client(), nil)
		server.srv.Handler = server.Handler(usm)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"email": "oscar@newemail.com"}`))
		req.Header.Add("Authorization", "Bearer 084a4a0f-92cc-46e6-9b57-1d2aed3c389e")
		server.srv.Handler.ServeHTTP(w, req)

		body, _ := ioutil.ReadAll(w.Body)
		err, _ := rest_errors.NewRestErrorFromBytes(body)
		// var restErr rest_errors.RestErr
		// json.Unmarshal(body, &restErr)
		//
		assert.NotNil(t, err)
	})
}

// correct... mocks a succesfull call to the oauth-api made inside of the middleware
func correctTestServerResponse() *httptest.Server {
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"access_token":"084a4a0f-92cc-46e6-9b57-1d2aed3c389e", "user_id":1, "user_role":"user", "expires":1637510344}`))
	}))
	testServer.Listener.Close()
	l, _ := net.Listen("tcp", "127.0.0.1:8081")
	testServer.Listener = l
	testServer.Start()
	return testServer
}
