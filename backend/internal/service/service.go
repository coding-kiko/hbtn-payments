package service

import (
	"control-pago-backend/internal/email"
	"control-pago-backend/internal/entity"
	"control-pago-backend/internal/repository"
	"control-pago-backend/log"
)

type Service interface {
	RegisterPayment(req *entity.RegisterPaymentRequest) error
}

type service struct {
	Repo        repository.Repository
	logger      log.Logger
	EmailClient email.EmailClient
}

func NewService(emailClient email.EmailClient, lgr log.Logger, repo repository.Repository) Service {
	return &service{
		Repo:        repo,
		EmailClient: emailClient,
		logger:      lgr,
	}
}
