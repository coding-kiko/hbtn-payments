package entity

type GetSummaryResponse struct {
	Summary string `json:"summary"`
}

type RegisterPaymentRequest struct {
	Month         string  `json:"month"`
	Amount        int     `json:"amount"`
	ReceiptBase64 string  `json:"receipt"`
	Company       string  `json:"company"`
	EmailTo       *string `json:"emailto,omitempty"`
}

type RegisterPayment struct {
	Month   string
	Amount  int
	Receipt string
	Company string
}

type Payment struct {
	Month   string
	Amount  int
	Receipt string
	Company string
}

type Receipt struct {
	Name string
	Data []byte
}
