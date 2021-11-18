package repositories

// TODO: pulish implementation, add proper error handling and finish testing

import (
	"database/sql"
	"strings"
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
	db  *sql.DB
	log ports.UserLogger
}

func NewUsersRepository(db *sql.DB, logger ports.UserLogger) ports.UsersRepository {
	onceUsersRepo.Do(func() {
		instanceUsersRepo = &usersRepository{
			db:  db,
			log: logger,
		}
	})
	return instanceUsersRepo
}

const (
	queryGetUser         = "SELECT id, first_name, last_name, email, date_created, status, role FROM users WHERE id=?;"
	queryGetUserByEmail  = "SELECT id, first_name, last_name, email, date_created, status, role, password FROM users WHERE email=?;"
	queryInsertUser      = "INSERT INTO users(first_name, last_name, email, date_created, status, password, role) VALUES(?, ?, ?, ?, ?, ?, ?);"
	queryUpdateUser      = "UPDATE users SET first_name=?, last_name=?, email=?, password=?, last_modified=? WHERE id=?;"
	queryUpdateUserAdmin = "UPDATE users SET first_name=?, last_name=?, email=?, password=?, status=?, role=?, last_modified=? WHERE id=?;"
	queryDeleteUser      = "UPDATE users SET status='inactive' WHERE id=?;"
)

const (
	errNoRow = "no rows in result"
)

func (r *usersRepository) Get(id int64) (*domain.User, rest_errors.RestErr) {
	stmt, err := r.db.Prepare(queryGetUser)
	if err != nil {
		r.log.Error(err.Error(), err)
		return nil, rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	var user domain.User
	result := stmt.QueryRow(id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Role); err != nil {
		r.log.Error(err.Error(), err)
		if strings.Contains(err.Error(), errNoRow) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		return nil, rest_errors.NewInternalServerError("db error")
	}

	return &user, nil
}

func (r *usersRepository) GetByEmail(email string) (*domain.User, rest_errors.RestErr) {
	stmt, err := r.db.Prepare(queryGetUserByEmail)
	if err != nil {
		r.log.Error(err.Error(), err)
		return nil, rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	var user domain.User
	result := stmt.QueryRow(email)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Role, &user.Password); err != nil {
		r.log.Error(err.Error(), err)
		if strings.Contains(err.Error(), errNoRow) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		return nil, rest_errors.NewInternalServerError("db error")
	}
	return &user, nil
}

func (r *usersRepository) Save(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryInsertUser)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password, user.Role)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	user.Id = userId
	return nil
}

func (r *usersRepository) Update(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryUpdateUser)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.LastModified, user.Id)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}

	return nil
}

func (r *usersRepository) UpdateAdmin(user *domain.User) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryUpdateUserAdmin)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Status, user.Role, user.LastModified, user.Id)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}

	return nil
}

func (r *usersRepository) Delete(id int64) rest_errors.RestErr {
	stmt, err := r.db.Prepare(queryDeleteUser)
	if err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id); err != nil {
		r.log.Error(err.Error(), err)
		return rest_errors.NewInternalServerError("db error")
	}
	return nil
}
