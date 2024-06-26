package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Repository repository.RepositoryInterface
	JWT        JWT
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	JWT        JWT
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		JWT:        opts.JWT,
	}
}

// RegisterHandlers to register all endpoints
func (server *Server) RegisterHandlers(e *echo.Echo) {
	e.POST("/register", func(c echo.Context) error {
		return server.RegisterUser(c)
	})

	e.POST("/login", func(c echo.Context) error {
		return server.LoginUser(c)
	})

	e.GET("/profile", func(c echo.Context) error {
		return server.GetUserProfile(c)
	})

	e.PUT("/profile", func(c echo.Context) error {
		return server.UpdateUserProfile(c)
	})
}
