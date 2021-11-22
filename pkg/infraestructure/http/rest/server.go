package rest

import (
	"database/sql"
	"net/http"

	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"github.com/FacuBar/bookstore_users-api/pkg/core/service"
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/repositories"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *sql.DB
	l      ports.UserLogger
	rest   *http.Client
	router *gin.Engine
}

func NewServer(db *sql.DB, l ports.UserLogger, rest *http.Client) *Server {
	server := &Server{
		db:   db,
		l:    l,
		rest: rest,
	}

	ur := repositories.NewUsersRepository(db, l)
	us := service.NewUsersService(ur)

	router := server.Handler(us)

	server.router = router
	return server
}

func (s *Server) Start(address string) {
	s.router.Run(address)
}
