package ports

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
)

type UsersService interface {
	GetUser(int64) (*domain.User, rest_errors.RestErr)
	Register(*domain.User) rest_errors.RestErr
	Update(*domain.User, bool) rest_errors.RestErr
	Login(string, string) (*domain.User, rest_errors.RestErr)

	// Logout() *rest_errors.RestErr
}
