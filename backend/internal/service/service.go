package service

import (
	"control-pago-backend/internal/email"
	"control-pago-backend/internal/entity"
	"control-pago-backend/internal/repository"
	"control-pago-backend/log"
)

type Service interface {
	RegisterPayment(req *entity.RegisterPaymentRequest) error
	GetSummary() (*entity.GetSummaryResponse, error)
}

type service struct {
	logger              log.Logger
	Repo                repository.Repository
	EmailClient         email.EmailClient
	StaticServerBaseUrl string
}

func NewService(emailClient email.EmailClient, lgr log.Logger, repo repository.Repository, staticServerBaseUrl string) Service {
	return &service{
		Repo:                repo,
		EmailClient:         emailClient,
		logger:              lgr,
		StaticServerBaseUrl: staticServerBaseUrl,
	}
}
