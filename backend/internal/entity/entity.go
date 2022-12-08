package entity

type GetSummaryResponse struct {
	Summary string `json:"summary"`
}

type RegisterPaymentRequest struct {
	Month   string  `json:"month"`
	Amount  int     `json:"amount"`
	Receipt string  `json:"receipt"`
	Company *string `json:"company,omitempty"`
	EmailTo string  `json:"emailto"`
}

type Receipt struct {
	Name string
	Data []byte
}
