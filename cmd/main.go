package main

import (
	"os"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/handler/helper/token"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	prvKey, err := os.ReadFile("id_rsa")
	if err != nil {
		e.Logger.Fatal(err)
	}
	pubKey, err := os.ReadFile("id_rsa.pub")
	if err != nil {
		e.Logger.Fatal(err)
	}

	opts := handler.NewServerOptions{
		Repository: repo,
		JWT:        token.NewJWT(prvKey, pubKey),
	}

	handler.NewServer(opts).RegisterHandlers(e)

	e.Logger.Fatal(e.Start(":1323"))
}
