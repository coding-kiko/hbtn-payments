package service

import "control-pago-backend/internal/entity"

func (s *service) GetSummary() (*entity.GetSummaryResponse, error) {
	payments, err := s.Repo.GetPayments()
	if err != nil {
		s.logger.Error("get_summay.go", "GetSummary", err.Error())
		return nil, err
	}

	summaryBase64, err := GenerateHtmlSummaryBase64(payments)
	if err != nil {
		s.logger.Error("get_summay.go", "GetSummary", err.Error())
		return nil, err
	}

	res := &entity.GetSummaryResponse{
		Summary: summaryBase64,
	}

	return res, nil
}

func GenerateHtmlSummaryBase64([]entity.Payment) (string, error) {
	return "", nil
}
