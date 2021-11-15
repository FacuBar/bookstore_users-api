package ports

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
)

type UsersService interface {
	Register(*domain.User) rest_errors.RestErr
	Update(*domain.User) rest_errors.RestErr
	Login(string, string) (*domain.User, rest_errors.RestErr)

	// Logout() *rest_errors.RestErr
}
