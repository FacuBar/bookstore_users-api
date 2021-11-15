package ports

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
)

type UsersRepository interface {
	Get(int64) (*domain.User, rest_errors.RestErr)
	Delete(int64) rest_errors.RestErr
	GetByEmail(string) (*domain.User, rest_errors.RestErr)

	Save(*domain.User) rest_errors.RestErr
	Update(*domain.User) rest_errors.RestErr
	UpdateAdmin(*domain.User) rest_errors.RestErr
}
