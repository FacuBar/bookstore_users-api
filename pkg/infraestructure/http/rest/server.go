package rest

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_users-api/pkg/core/service"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/repositories"
)

type Server struct {
	db   *sql.DB
	l    ports.UserLogger
	rest *http.Client
	srv  *http.Server
}

func NewServer(srv *http.Server, db *sql.DB, l ports.UserLogger, rest *http.Client) *Server {
	server := &Server{
		db:   db,
		l:    l,
		rest: rest,
		srv:  srv,
	}

	ur := repositories.NewUsersRepository(db, l)
	us := service.NewUsersService(ur)

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

	go func() {
		if err := s.srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

	log.Println("Server exiting")
}
