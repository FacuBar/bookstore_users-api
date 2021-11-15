package mysql

// TODO: pulish implementation, add proper error handling and finish testing

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

const (
	queryGetUser         = "SELECT id, first_name, last_name, email, date_created, status, privileges FROM users WHERE id=?;"
	queryGetUserByEmail  = "SELECT id, first_name, last_name, email, date_created, status, privileges, password FROM users WHERE email=?;"
	queryInsertUser      = "INSERT INTO users(first_name, last_name, email, date_created, status, password, privileges) VALUES(?, ?, ?, ?, ?, ?, ?);"
	queryUpdateUser      = "UPDATE users SET first_name=?, last_name=?, email=?, password=? WHERE id=?;"
	queryUpdateUserAdmin = "UPDATE users SET first_name=?, last_name=?, email=?, password=?, status=?, privilege=? WHERE id=?;"
	queryDeleteUser      = "DELETE FROM users WHERE id=?;"
)

func (r *usersRepository) Get(id int64) (*domain.User, rest_errors.RestErr) {
	stmt, err := r.db.Prepare(queryGetUser)
	if err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	var user domain.User
	result := stmt.QueryRow(id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Privileges); err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error())
	}

	return &user, nil
}

func (r *usersRepository) Save(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryInsertUser)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password, user.Privileges)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	user.Id = userId
	return nil
}

func (r *usersRepository) Update(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryUpdateUser)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Id)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *usersRepository) UpdateAdmin(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryUpdateUser)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Status, user.Privileges, user.Id)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *usersRepository) Delete(id int64) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryDeleteUser)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id); err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *usersRepository) GetByEmail(email string) (*domain.User, rest_errors.RestErr) {
	stmt, err := r.db.Prepare(queryGetUserByEmail)
	if err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	var user domain.User
	result := stmt.QueryRow(email)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Privileges, &user.Password); err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}
