package ports

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
)

type UsersService interface {
	Register(domain.User) *domain.User
	Login(domain.User) *domain.User
	Update(domain.User) *domain.User

	// Logout() *rest_errors.RestErr
}
