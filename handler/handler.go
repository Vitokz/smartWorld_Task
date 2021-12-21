package handler

import (
	"github.com/Vitokz/smartWorld_Task/config"
	"github.com/Vitokz/smartWorld_Task/internal/repository"
	"github.com/Vitokz/smartWorld_Task/internal/services/jwt"
	"github.com/Vitokz/smartWorld_Task/internal/services/logger"
	"net/http"
	"time"
)

type Handler struct {
	log    logger.Logger
	Config *config.Config
	Client HTTPClient

	Repository repository.Repository
	Jwt        jwt.JwtService
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHandler(cfg *config.Config, logger logger.Logger, repo repository.Repository) *Handler {
	return &Handler{
		log:    logger,
		Config: cfg,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Repository: repo,
		Jwt:        jwt.NewJwtService(cfg),
	}
}
