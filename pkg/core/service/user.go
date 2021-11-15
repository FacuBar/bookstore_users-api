package service

import (
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
}

func NewUsersService(repo ports.UsersRepository) ports.UsersService {
	onceUsersService.Do(func() {
		instanceUsersService = &usersService{
			repo: repo,
		}
	})
	return instanceUsersService
}

var (
	statusActive      = "active"
	defaultPriveleges = 1 << 0
	dateLayout        = "2006-01-02 15:04:05"
)

func (s *usersService) Register(user *domain.User) rest_errors.RestErr {
	if err := Validate(user); err != nil {
		return err
	}

	user.DateCreated = time.Now().UTC().Format(dateLayout)
	user.Status = statusActive
	user.Privileges = defaultPriveleges

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)

	if err := s.repo.Save(user); err != nil {
		return err
	}
	return nil
}

func (s *usersService) Login(email string, password string) (*domain.User, rest_errors.RestErr) {
	user, err := s.repo.GetByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, rest_errors.NewBadRequestError("Invalid credentials")
	}
	return user, nil
}

func (s *usersService) Update(user *domain.User) rest_errors.RestErr {
	oldUser, err := s.repo.Get(user.Id)
	if err != nil {
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
	if user.Password == "" {
		user.Password = oldUser.Password
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hashedPassword)
	}

	user.LastModified = time.Now().UTC().Format(dateLayout)

	if err := s.repo.Update(user); err != nil {
		return err
	}
	return nil
}

// func (s *usersService) Logout() *rest_errors.RestErr {
// 	return nil
// }

func Validate(u *domain.User) rest_errors.RestErr {
	Format(u)
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if u.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}

func Format(u *domain.User) {
	u.Email = strings.ToLower((strings.TrimSpace(u.Email)))
	u.FirstName = strings.ToTitle(strings.TrimSpace(u.FirstName))
	u.LastName = strings.ToTitle((strings.TrimSpace(u.LastName)))
}
