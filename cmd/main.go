package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/clients"
	oauth_grpc "github.com/FacuBar/bookstore_users-api/pkg/infraestructure/http/grpc/oauth"
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

	oauthClient, err := oauth_grpc.NewClient("0.0.0.0:10000")
	if err != nil {
		panic("error initializing grpc client")
	}
	fmt.Print("testing block")

	server := rest.NewServer(&http.Server{Addr: ":8080"}, db, l, &http.Client{}, oauthClient)

	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Stop(ctx)
}
