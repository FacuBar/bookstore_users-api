package rest

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
)

// swagger:response genericUser
type genericUserResponse struct {
	// in: body
	Body domain.User
}

// swagger:response genericError
type genericErrorResponse struct {
	// in: body
	Body rest_errors.RestErr
}

// swagger:parameters listUser
type userIDParamsWrapper struct {
	// in: path
	// required: true
	// minimum : 1
	ID int `json:"id"`
}

// swagger:parameters listUser
type authHeaderWrapper struct {
	// in: header
	// example: "Bearer {auth_token}"
	Authorization string `json:"Authorization"`
}

type requestRegister struct {
	// required : true
	// example : Oscar
	FirstName string `json:"first_name"`
	// required : true
	// example : isaac
	LastName string `json:"last_name"`
	// required : true
	// example : oscaac@email.com
	Email string `json:"email"`
	// required : true
	// example : somepassword
	Password string `json:"password"`
	// required : true
	// example : somepassword
	ConfirmPassword string `json:"confirm_password"`
}

// swagger:parameters registerUsers
type requestRegisterWrapper struct {
	// in: body
	Body requestRegister
}

type requestLoginUser struct {
	// example : user1@email.com
	// required : true
	Email string `json:"email"`
	// example : somepassword
	// required : true
	Password string `json:"password"`
}

// swagger:parameters loginUsers
type requestLoginWrapper struct {
	// in: body
	Body requestLoginUser
}
