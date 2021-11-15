package mysql

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/stretchr/testify/assert"
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
	Privileges:  1,
}

func TestGet(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db}
		defer func() {
			repo.db.Close()
		}()

		query := "SELECT id, first_name, last_name, email, date_created, status, privileges FROM users WHERE id=\\?;"

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "date_created", "status", "privileges"}).
			AddRow(test.Id, test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Privileges)

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(test.Id).WillReturnRows(rows)

		user, err := repo.Get(test.Id)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, test.FirstName, user.FirstName)
		assert.EqualValues(t, test.Email, user.Email)
	})
}

func TestSave(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db}
		defer func() {
			repo.db.Close()
		}()

		query := "INSERT INTO users\\(first_name, last_name, email, date_created, status, password, privileges\\) VALUES\\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\);"

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, test.Email, test.DateCreated, test.Status, test.Password, test.Privileges).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(&test)

		assert.Nil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		db, mock := NewMock()

		repo := &usersRepository{db: db}
		defer func() {
			repo.db.Close()
		}()

		query := "UPDATE users SET first_name=\\?, last_name=\\?, email=\\?, password=\\? WHERE id=\\?;"

		mock.ExpectPrepare(query).ExpectExec().WithArgs(test.FirstName, test.LastName, "random@gmail.com", test.Password, test.Id).WillReturnResult(sqlmock.NewResult(1, 1))

		test.Email = "random@gmail.com"

		err := repo.Update(&test)

		assert.Nil(t, err)
	})
}
