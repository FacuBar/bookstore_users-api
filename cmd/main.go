package main

import (
	"github.com/FacuBar/bookstore_users-api/pkg/core/service"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/clients"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/http/rest"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/logger"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/repositories"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db := clients.ConnectDB()
	l := logger.NewUserLogger()
	ur := repositories.NewUsersRepository(db, l)
	us := service.NewUsersService(ur)

	router := rest.Handler(us)
	router.Run(":8080")
}
