package service

import (
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/FacuBar/bookstore_users-api/pkg/core/domain"
	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_utils-go/rest_errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	onceUsersService     sync.Once
	instanceUsersService *usersService
)

type usersService struct {
	repo ports.UsersRepository
	rmq  ports.UserRMQ
}

func NewUsersService(repo ports.UsersRepository, rmq ports.UserRMQ) ports.UsersService {
	onceUsersService.Do(func() {
		instanceUsersService = &usersService{
			repo: repo,
			rmq:  rmq,
		}
	})
	return instanceUsersService
}

var (
	statusActive = "active"
	defaultRole  = "user"
	dateLayout   = "2006-01-02 15:04:05"
)

func (s *usersService) GetUser(userId int64) (*domain.User, rest_errors.RestErr) {
	user, err := s.repo.Get(userId)
	if err != nil {
		if err.Status() == http.StatusInternalServerError {
			return nil, rest_errors.NewInternalServerError("error while trying to get user, try again later")
		}
		return nil, err
	}

	return user, nil
}

func (s *usersService) Register(user *domain.User) rest_errors.RestErr {
	if err := validate(user); err != nil {
		return err
	}

	user.DateCreated = time.Now().UTC().Format(dateLayout)
	user.Status = statusActive
	user.Role = defaultRole

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)

	if err := s.repo.Save(user); err != nil {
		if err.Status() != http.StatusInternalServerError {
			return err
		}
		return rest_errors.NewInternalServerError("error while trying to register, try again later")
	}
	s.rmq.Publish("users.event.register", user)
	return nil
}

func (s *usersService) Login(email string, password string) (*domain.User, rest_errors.RestErr) {
	user, err := s.repo.GetByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		switch err.Status() {
		case http.StatusInternalServerError:
			return nil, rest_errors.NewInternalServerError("error while trying to login, try again later")
		case http.StatusNotFound:
			return nil, rest_errors.NewBadRequestError("invalid credentials")
		default:
			return nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, rest_errors.NewBadRequestError("invalid credentials")
	}
	return user, nil
}

func (s *usersService) Update(user *domain.User, isAdmin bool) rest_errors.RestErr {
	oldUser, err := s.repo.Get(user.Id)
	if err != nil {
		if err.Status() == http.StatusInternalServerError {
			return rest_errors.NewInternalServerError("error while trying to fetch user, try again later")
		}
		return err
	}

	if strings.TrimSpace(user.FirstName) == "" {
		user.FirstName = oldUser.FirstName
	}
	if strings.TrimSpace(user.LastName) == "" {
		user.LastName = oldUser.LastName
	}
	if strings.TrimSpace(user.Email) == "" {
		user.Email = oldUser.Email
	}
	if user.Status == "" {
		user.Status = oldUser.Status
	}
	if user.Role == "" {
		user.Role = oldUser.Role
	}
	if user.Password == "" {
		user.Password = oldUser.Password
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hashedPassword)
	}

	user.LastModified = time.Now().UTC().Format(dateLayout)

	if isAdmin {
		if err := s.repo.UpdateAdmin(user); err != nil {
			return rest_errors.NewInternalServerError("error while trying to update user, try again later")
		}
	}
	if err := s.repo.Update(user); err != nil {
		return rest_errors.NewInternalServerError("error while trying to update user, try again later")
	}
	s.rmq.Publish("users.event.update", user)
	return nil
}

// func (s *usersService) Logout() *rest_errors.RestErr {
// 	return nil
// }

func validate(u *domain.User) rest_errors.RestErr {
	format(u)
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if u.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	if u.FirstName == "" {
		return rest_errors.NewBadRequestError("invalid first name")
	}
	if u.LastName == "" {
		return rest_errors.NewBadRequestError("invalid last name")
	}
	return nil
}

func format(u *domain.User) {
	u.Email = strings.ToLower((strings.TrimSpace(u.Email)))
	u.FirstName = strings.Title(strings.TrimSpace(u.FirstName))
	u.LastName = strings.Title((strings.TrimSpace(u.LastName)))
}
