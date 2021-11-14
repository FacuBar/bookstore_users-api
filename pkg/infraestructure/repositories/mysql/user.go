package mysql

import (
	"database/sql"
	"sync"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
)

var (
	onceUsersRepo     sync.Once
	instanceUsersRepo *usersRepository
)

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) ports.UsersRepository {
	onceUsersRepo.Do(func() {
		instanceUsersRepo = &usersRepository{
			db: db,
		}
	})
	return instanceUsersRepo
}

func (r *usersRepository) Get(id int64) (*domain.User, rest_errors.RestErr) {
	return nil, nil
}

func (r *usersRepository) Save(user *domain.User) rest_errors.RestErr {
	return nil
}

func (r *usersRepository) Update(user *domain.User) rest_errors.RestErr {
	return nil
}

func (r *usersRepository) UpdateAdmin(user *domain.User) rest_errors.RestErr {
	return nil
}

func (r *usersRepository) Delete(id int64) rest_errors.RestErr {
	return nil
}

func (r *usersRepository) GetByEmail(email string) (*domain.User, rest_errors.RestErr) {
	return nil, nil
}
