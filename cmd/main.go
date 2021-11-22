package main

import (
	"net/http"

	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/clients"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/http/rest"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db := clients.ConnectDB()
	l := logger.NewUserLogger()

	server := rest.NewServer(db, l, &http.Client{})

	server.Start(":8080")
}
