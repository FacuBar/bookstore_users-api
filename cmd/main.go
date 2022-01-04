package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/clients"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/http/rest"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/logger"
	"github.com/FacuBar/bookstore_utils-go/auth"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db := clients.ConnectDB()
	l := logger.NewUserLogger()

	oauthClient, err := auth.NewClient(os.Getenv("GRPC_OAUTH_ADDRESS"))
	if err != nil {
		panic("error connecting grpc client")
	}

	RMQ, err := clients.NewRabbitMQ(os.Getenv("RMQ_URI"))
	if err != nil {
		panic("error connecting rabbitMQ")
	}

	server := rest.NewServer(&http.Server{Addr: os.Getenv("PORT")}, db, l, oauthClient, RMQ)

	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Stop(ctx)
}
