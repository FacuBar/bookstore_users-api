package service

import (
	"net/http"
	"testing"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type userRepoMock struct{}

var (
	funcGet         func(int64) (*domain.User, rest_errors.RestErr)
	funcDelete      func(int64) rest_errors.RestErr
	funcGetByEmail  func(string) (*domain.User, rest_errors.RestErr)
	funcSave        func(*domain.User) rest_errors.RestErr
	funcUpdate      func(*domain.User) rest_errors.RestErr
	funcUpdateAdmin func(*domain.User) rest_errors.RestErr
)

var (
	repoMock ports.UsersRepository = &userRepoMock{}
)

func (m *userRepoMock) Get(id int64) (*domain.User, rest_errors.RestErr) {
	return funcGet(id)
}
func (m *userRepoMock) Delete(id int64) rest_errors.RestErr {
	return funcDelete(id)
}
func (m *userRepoMock) GetByEmail(email string) (*domain.User, rest_errors.RestErr) {
	return funcGetByEmail(email)
}
func (m *userRepoMock) Save(user *domain.User) rest_errors.RestErr {
	return funcSave(user)
}
func (m *userRepoMock) Update(user *domain.User) rest_errors.RestErr {
	return funcUpdate(user)
}
func (m *userRepoMock) UpdateAdmin(user *domain.User) rest_errors.RestErr {
	return funcUpdateAdmin(user)
}

var (
	userTest = domain.User{
		Id:          1,
		FirstName:   "Oscar",
		LastName:    "Isaac",
		Email:       "oscaac@gmail.com",
		Password:    "$2a$10$jRL.gYiodDnwcOBErnDfuu5044h40PM3ZOAOzit6O4RIL9wG24xJ6", //password
		DateCreated: "2006-01-02 15:04:05",
		Status:      "active",
		Role:        "user",
	}
)

func TestGetUser(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}

		s := NewUsersService(repoMock)
		user, err := s.GetUser(1)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, "Oscar", user.FirstName)
		assert.EqualValues(t, "oscaac@gmail.com", user.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}

		s := NewUsersService(repoMock)
		user, err := s.GetUser(1)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "user not found", err.Message())
	})

	t.Run("DbError", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("db error")
		}

		s := NewUsersService(repoMock)
		user, err := s.GetUser(1)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "error while trying to get user, try again later", err.Message())
	})
}

func TestRegister(t *testing.T) {
	t.Run("UserNotValid", func(t *testing.T) {
		s := NewUsersService(repoMock)

		user := domain.User{Email: ""}
		err := s.Register(&user)

		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid email address", err.Message())
	})

	t.Run("NoError", func(t *testing.T) {
		funcSave = func(u *domain.User) rest_errors.RestErr {
			return nil
		}

		s := NewUsersService(repoMock)

		user := userTest
		user.Password = "password"
		err := s.Register(&user)

		assert.Nil(t, err)
		errB := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
		assert.Nil(t, errB)
	})

	t.Run("DbError", func(t *testing.T) {
		funcSave = func(u *domain.User) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("db error")
		}
		s := NewUsersService(repoMock)

		user := userTest
		err := s.Register(&user)

		assert.NotNil(t, err)
		assert.EqualValues(t, "error while trying to register, try again later", err.Message())
	})
}

func TestLogin(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		funcGetByEmail = func(s string) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}
		s := NewUsersService(repoMock)

		user, err := s.Login("oscaac@gmail.com", "password")

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, "Oscar", user.FirstName)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		funcGetByEmail = func(s string) (*domain.User, rest_errors.RestErr) {
			return &userTest, nil
		}
		s := NewUsersService(repoMock)

		user, err := s.Login("oscaac@gmail.com", "notthepassword")

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid credentials", err.Message())
		assert.EqualValues(t, http.StatusBadRequest, err.Status())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		funcGetByEmail = func(s string) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		s := NewUsersService(repoMock)

		user, err := s.Login("oscaac@gmail.com", "password")

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid credentials", err.Message())
		assert.EqualValues(t, http.StatusBadRequest, err.Status())
	})

	t.Run("DbError", func(t *testing.T) {
		funcGetByEmail = func(s string) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("db error")
		}
		s := NewUsersService(repoMock)

		user, err := s.Login("oscaac@gmail.com", "password")

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, "error while trying to login, try again later", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestUpdate(t *testing.T) {
	userTestUpdate := domain.User{
		Id:        1,
		FirstName: "Llewyn",
		LastName:  "Davis",
		Password:  "Ulysses",
	}

	t.Run("NoError", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			user := userTest
			return &user, nil
		}

		funcUpdate = func(u *domain.User) rest_errors.RestErr {
			return nil
		}

		s := NewUsersService(repoMock)

		update := userTestUpdate

		err := s.Update(&update)

		assert.Nil(t, err)
		assert.EqualValues(t, "oscaac@gmail.com", update.Email)
		assert.EqualValues(t, "Llewyn", update.FirstName)

		errB := bcrypt.CompareHashAndPassword([]byte(update.Password), []byte("Ulysses"))
		assert.Nil(t, errB)
	})

	t.Run("SearchError", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("db error")
		}

		s := NewUsersService(repoMock)

		update := userTestUpdate

		err := s.Update(&update)

		assert.NotNil(t, err)
		assert.EqualValues(t, "error while trying to fetch user, try again later", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}

		s := NewUsersService(repoMock)

		update := userTestUpdate

		err := s.Update(&update)

		assert.NotNil(t, err)
		assert.EqualValues(t, "user not found", err.Message())
		assert.EqualValues(t, http.StatusNotFound, err.Status())
	})

	t.Run("ErrorSave", func(t *testing.T) {
		funcGet = func(i int64) (*domain.User, rest_errors.RestErr) {
			user := userTest
			return &user, nil
		}

		funcUpdate = func(u *domain.User) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("db error")
		}
		s := NewUsersService(repoMock)

		update := userTestUpdate

		err := s.Update(&update)

		assert.NotNil(t, err)
		assert.EqualValues(t, "error while trying to update user, try again later", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	})
}

func TestValidate(t *testing.T) {
	t.Run("InvalidEmail", func(t *testing.T) {
		testEm := domain.User{Email: "asd@adasd,d", Password: "a", FirstName: "b", LastName: "c"}
		err := validate(&testEm)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid email address", err.Message())
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		testPass := domain.User{Email: "asd@gmail.com", Password: "", FirstName: "b", LastName: "c"}
		err := validate(&testPass)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid password", err.Message())
	})

	t.Run("InvalidFirstName", func(t *testing.T) {
		testFn := domain.User{Email: "asd@gmail.com", Password: "a", FirstName: "  ", LastName: "c"}
		err := validate(&testFn)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid first name", err.Message())
	})

	t.Run("InvalidLastName", func(t *testing.T) {
		testLn := domain.User{Email: "asd@gmail.com", Password: "a", FirstName: "b", LastName: " "}
		err := validate(&testLn)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid last name", err.Message())
	})
}
