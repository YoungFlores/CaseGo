package server

import (
	"net/http"

	"github.com/YoungFlores/Case_Go/Profile/internal/api"
	"github.com/YoungFlores/Case_Go/Profile/internal/db"
	profileHandler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers"
	profileRepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
	profileService "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
)

type Sever struct {
	HTTP *http.Server
	DB   *db.DataBase
}

func New() (*Sever, error) {

	database := &db.DataBase{}
	config := LoadConfig()

	if err := database.Open(
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	); err != nil {
		return nil, err
	}

	pr := profileRepo.NewPostgresProfileRepo(database.GetDB())

	ps := profileService.NewProfileService(pr)

	jwtMiddleware := rs256.New(config.PublicKey, "auth", "all")

	profileHandlers := profileHandler.NewProfileHandler(ps)

	router := api.SetupRouter(profileHandlers, jwtMiddleware)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &Sever{
		HTTP: srv,
		DB:   database,
	}, nil
}
