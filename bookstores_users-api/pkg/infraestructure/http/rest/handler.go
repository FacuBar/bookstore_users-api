// Package classifiction Bookstore-users API.
//
// Documentation for Bookstore-users API
//
// Schemes: http
// Host: localhost
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json

// swagger:meta
package rest

import (
	"net/http"
	"strconv"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) Handler(us ports.UsersService) *gin.Engine {
	router := gin.Default()

	router.POST("/users", registerUser(us))
	router.POST("/users/login", login(us))
	router.GET("/users/:user_id", authenticate(getUser(us), s.rest))
	router.PUT("/users/:user_id", authenticate(updateUser(us), s.rest))
	// router.DELETE("/users/:user_id")

	// Paymentoptions relatedendpoints ("/users/:user_id/paymentoptions...")

	// Serving docs
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	router.GET("/docs", gin.WrapH(sh))
	router.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./"))))

	return router
}

// swagger:route POST /users users registerUsers
// Registers a new user into the database
// responses:
// 	200: genericUser
// 	400: genericError
// 	500: genericError
func registerUser(s ports.UsersService) gin.HandlerFunc {
	type request struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	return func(c *gin.Context) {
		var registerRequest request
		if err := c.ShouldBindJSON(&registerRequest); err != nil {
			restErr := rest_errors.NewBadRequestError("invalid request")
			c.JSON(restErr.Status(), restErr)
			return
		}

		if registerRequest.Password != registerRequest.ConfirmPassword {
			restErr := rest_errors.NewBadRequestError("passwords are not equal")
			c.JSON(restErr.Status(), restErr)
			return
		}

		user := domain.User{
			FirstName: registerRequest.FirstName,
			LastName:  registerRequest.LastName,
			Email:     registerRequest.Email,
			Password:  registerRequest.Password,
		}

		if err := s.Register(&user); err != nil {
			c.JSON(err.Status(), err)
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

// swagger:route GET /users/{user_id} users listUser
// List information of a particular user
// Only accessible by the authenticated user
// responses:
// 	200: genericUser
// 	400: genericError
// 	401: genericError
// 	404: genericError
// 	500: genericError
func getUser(s ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if userErr != nil {
			restErr := rest_errors.NewBadRequestError("user id not valid")
			c.JSON(restErr.Status(), restErr)
			return
		}

		authorizedUser := c.MustGet(userPayloadKey).(userPayload)
		if authorizedUser.Id != userId {
			restErr := rest_errors.NewUnauthorizedError("you don't have the permissions to access this resource")
			c.JSON(restErr.Status(), restErr)
			return
		}

		user, err := s.GetUser(userId)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// swagger:route POST /users/login users loginUsers
// Validates that the email and the passwords provided are valid for a registered user
// responses:
// 	200: genericUser
// 	400: genericError
// 	500: genericError
func login(s ports.UsersService) gin.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var loginRequest request
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			restErr := rest_errors.NewBadRequestError("invalid request")
			c.JSON(restErr.Status(), restErr)
			return
		}

		user, err := s.Login(loginRequest.Email, loginRequest.Password)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// swagger:route PUT /users/{user_id} users updateUser
// Validates that the email and the passwords provided are valid for a registered user
// responses:
// 	200: genericUser
// 	400: genericError
//  401: genericError
// 	500: genericError
func updateUser(s ports.UsersService) gin.HandlerFunc {
	type request struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		Status          string `json:"status"`
		Role            string `json:"role"`
	}

	return func(c *gin.Context) {
		userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if userErr != nil {
			restErr := rest_errors.NewBadRequestError("user id not valid")
			c.JSON(restErr.Status(), restErr)
			return
		}

		authorizedUser := c.MustGet(userPayloadKey).(userPayload)

		if authorizedUser.Id != userId && authorizedUser.Role != "admin" {
			restErr := rest_errors.NewUnauthorizedError("you don't have the permissions to perform this action")
			c.JSON(restErr.Status(), restErr)
			return
		}

		var userRequest request
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			restErr := rest_errors.NewBadRequestError("invalid request body")
			c.JSON(restErr.Status(), restErr)
			return
		}

		if userRequest.Password != userRequest.ConfirmPassword {
			restErr := rest_errors.NewBadRequestError("passwords are not equal")
			c.JSON(restErr.Status(), restErr)
			return
		}

		user := domain.User{
			Id:        userId,
			FirstName: userRequest.FirstName,
			LastName:  userRequest.LastName,
			Email:     userRequest.Email,
			Password:  userRequest.Password,
		}

		if authorizedUser.Role == "admin" {
			user.Role = userRequest.Role
			user.Status = userRequest.Status
		}

		if err := s.Update(&user, authorizedUser.Role == "admin"); err != nil {
			c.JSON(err.Status(), err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
