package repositories

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var test = domain.User{
	Id:          1,
	FirstName:   "Oscar",
	LastName:    "Isaac",
	Email:       "oscaac@gmail.com",
	Password:    "x4BRvJE8glEHeAX8GkevhxTNsglMxpIBdSjXj4O6538jdqbzx0saVbfJc3hZ",
	DateCreated: "2006-01-02 15:04:05",
	Status:      "active",
	Role:        "user",
}

type loggerMock struct {
}

func (l *loggerMock) Error(msg string, err error, tags ...zap.Field) {}
func (l *loggerMock) Info(msg string, tags ...zap.Field)             {}

func TestGet(t *testing.T) {
	query := "SELECT id, first_name, last_name, email, date_created, status, role FROM users WHERE id=\\?;"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "date_created", "status", "role"}).
			AddRow(test.Id, test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Role)

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Id).WillReturnRows(rows)

		user, err := repo.Get(test.Id)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, test.FirstName, user.FirstName)
		assert.EqualValues(t, test.Email, user.Email)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Id).WillReturnError(errors.New("error: no rows in result set"))

		user, err := repo.Get(test.Id)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "user not found", err.Message())
		assert.EqualValues(t, http.StatusNotFound, err.Status())
	})

	t.Run("QueryingError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Id).WillReturnError(errors.New("error"))
		user, err := repo.Get(test.Id)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("PrepareError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		user, err := repo.Get(test.Id)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestGetByEmail(t *testing.T) {
	query := "SELECT id, first_name, last_name, email, date_created, status, role, password FROM users WHERE email=\\?;"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "date_created", "status", "role", "password"}).
			AddRow(test.Id, test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Role, test.Password)
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Email).WillReturnRows(rows)

		user, err := repo.GetByEmail(test.Email)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, test.FirstName, user.FirstName)
		assert.EqualValues(t, test.Email, user.Email)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Email).WillReturnError(errors.New("error: no rows in result set"))

		user, err := repo.GetByEmail(test.Email)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "user not found", err.Message())
		assert.EqualValues(t, http.StatusNotFound, err.Status())
	})

	t.Run("QueryingError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Email).WillReturnError(errors.New("error"))
		user, err := repo.GetByEmail(test.Email)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("PrepareError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		user, err := repo.GetByEmail(test.Email)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestSave(t *testing.T) {
	query := "INSERT INTO users\\(first_name, last_name, email, date_created, status, password, role\\) VALUES\\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\);"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Password, test.Role).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(&test)

		assert.Nil(t, err)
	})

	t.Run("SavingUser", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Password, test.Role).WillReturnError(errors.New("..."))

		err := repo.Save(&test)

		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("PrepareError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		err := repo.Save(&test)

		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestUpdate(t *testing.T) {
	query := "UPDATE users SET first_name=\\?, last_name=\\?, email=\\?, password=\\?, last_modified=\\? WHERE id=\\?;"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, "random@gmail.com", test.Password, "2006-01-02 15:04:05", test.Id).WillReturnResult(sqlmock.NewResult(1, 1))

		test.Email = "random@gmail.com"
		test.LastModified = "2006-01-02 15:04:05"

		err := repo.Update(&test)

		assert.Nil(t, err)
	})

	t.Run("ExecError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, test.Email, test.Password, test.LastModified, test.Id).WillReturnError(sql.ErrConnDone)

		err := repo.Update(&test)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("PrepareError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		err := repo.Update(&test)

		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestUpdateAdmin(t *testing.T) {
	query := "UPDATE users SET first_name=\\?, last_name=\\?, email=\\?, password=\\?, status=\\?, role=\\?, last_modified=\\? WHERE id=\\?;"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, "random@gmail.com", test.Password, test.Status, test.Role, "2006-01-02 15:04:05", test.Id).WillReturnResult(sqlmock.NewResult(1, 1))

		test.Email = "random@gmail.com"
		test.LastModified = "2006-01-02 15:04:05"

		err := repo.UpdateAdmin(&test)

		assert.Nil(t, err)
	})

	t.Run("ExecError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, test.Email, test.Password, test.Status, test.Role, test.LastModified, test.Id).WillReturnError(sql.ErrConnDone)

		err := repo.UpdateAdmin(&test)
		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("PrepareError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		err := repo.UpdateAdmin(&test)

		assert.NotNil(t, err)
		assert.EqualValues(t, "db error", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestDelete(t *testing.T) {
	query := "UPDATE users SET status='inactive' WHERE id=\\?;"

	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.Id).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(test.Id)

		assert.Nil(t, err)
	})

	t.Run("ErrorPrepare", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

		err := repo.Delete(test.Id)

		assert.NotNil(t, err)
	})

	t.Run("ErrorExec", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db, log: &loggerMock{}}
		defer func() {
			repo.db.Close()
		}()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.Id).WillReturnError(sql.ErrConnDone)

		err := repo.Delete(test.Id)

		assert.NotNil(t, err)
	})
}
