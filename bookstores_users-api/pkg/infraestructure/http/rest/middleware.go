package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey  = "Authorization"
	userPayloadKey = "user_payload"
)

type userPayload struct {
	Id   int64  `json:"user_id"`
	Role string `json:"user_role"`
}

func authenticate(handler gin.HandlerFunc, restClient *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authHeaderKey)
		if len(authorizationHeader) == 0 {
			err := rest_errors.NewBadRequestError("no authorization header was provided")
			c.AbortWithStatusJSON(err.Status(), err)
			return
		}

		authFields := strings.Split(authorizationHeader, " ")
		if len(authFields) != 2 {
			err := rest_errors.NewBadRequestError("invalid authorization header format")
			c.AbortWithStatusJSON(err.Status(), err)
			return
		}

		if authFields[0] != "Bearer" {
			err := rest_errors.NewBadRequestError("authorization type not supported")
			c.AbortWithStatusJSON(err.Status(), err)
			return
		}

		resp, err := restClient.Get(fmt.Sprintf("http://localhost:8081/oauth/access_token/%s", authFields[1]))
		if err != nil {
			restErr := rest_errors.NewInternalServerError("couldn't verify session's validity")
			c.AbortWithStatusJSON(restErr.Status(), restErr)
			return
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			restErr := rest_errors.NewInternalServerError("couldn't verify session's validity")
			c.AbortWithStatusJSON(restErr.Status(), restErr)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			if resp.StatusCode == http.StatusInternalServerError {
				restErr := rest_errors.NewInternalServerError("couldn't verify session's validity")
				c.AbortWithStatusJSON(restErr.Status(), restErr)
				return
			}
			restErr := rest_errors.NewUnauthorizedError("you are not logged in")
			c.AbortWithStatusJSON(restErr.Status(), restErr)
			return
		}

		var user userPayload
		json.Unmarshal(bytes, &user)
		c.Set(userPayloadKey, user)

		handler(c)
	}
}
