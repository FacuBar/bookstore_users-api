package rest

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_users-api/pkg/core/service"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/clients"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/repositories"
	"github.com/FacuBar/bookstore_utils-go/auth"
)

type Server struct {
	db       *sql.DB
	l        ports.UserLogger
	srv      *http.Server
	oauth    *auth.Client
	rabbitmq *clients.RabbitMQ
}

func NewServer(srv *http.Server, db *sql.DB, l ports.UserLogger, oauth *auth.Client, rmq *clients.RabbitMQ) *Server {
	server := &Server{
		db:       db,
		l:        l,
		srv:      srv,
		oauth:    oauth,
		rabbitmq: rmq,
	}

	ur := repositories.NewUsersRepository(db, l)
	us := service.NewUsersService(ur, rmq)

	router := server.Handler(us)

	srv.Handler = router
	return server
}

func (s *Server) Start() {
	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) {
	s.db.Close()
	s.oauth.CC.Close()
	s.rabbitmq.Channel.Close()
	s.rabbitmq.Connection.Close()

	go func() {
		if err := s.srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

	log.Println("Server exiting")
	os.Exit(0)
}
