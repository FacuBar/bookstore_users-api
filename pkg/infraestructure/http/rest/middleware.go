package rest

// import (
// 	"strings"

// 	"github.com/FacuBar/bookstore_utils-go/rest_errors"
// 	"github.com/gin-gonic/gin"
// )

// const (
// 	authHeaderKey = "Authorization"
// )

// func authenticate(handler gin.HandlerFunc, restClient *http.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authorizationHeader := c.GetHeader(authHeaderKey)
// 		if len(authorizationHeader) == 0 {
// 			err := rest_errors.NewBadRequestError("no authorization header is provided")
// 			c.AbortWithStatusJSON(err.Status(), err)
// 		}

// 		authFields := strings.Split(authorizationHeader, " ")
// 		if len(authFields) != 2 {
// 			err := rest_errors.NewBadRequestError("invalid authorization header format")
// 			c.AbortWithStatusJSON(err.Status(), err)
// 		}

// 		if authFields[0] != "Bearer" {
// 			err := rest_errors.NewBadRequestError("authorization type not supported")
// 			c.AbortWithStatusJSON(err.Status(), err)
// 		}

// 	}
// }
