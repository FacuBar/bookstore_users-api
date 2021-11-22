package rest

import (
	"net/http"
	"strconv"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

func (s *Server) Handler(us ports.UsersService) *gin.Engine {
	router := gin.Default()

	router.POST("/users", registerUser(us))
	router.POST("/users/login", login(us))
	router.GET("/users/:user_id", getUser(us))

	return router
}

// User handlers
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
			restErr := rest_errors.NewInternalServerError("invalid request")
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

func getUser(s ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if userErr != nil {
			restErr := rest_errors.NewBadRequestError("user id not valid")
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
