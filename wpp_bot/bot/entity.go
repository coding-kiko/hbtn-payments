package bot

type GetPaymentSummaryResponse struct {
	SummaryBase64 string `json:"summary"`
}

type ReqBody struct {
	Month         string  `json:"month"`
	Amount        int     `json:"amount"`
	Email         string  `json:"emailto,omitempty"`
	Company       *string `json:"company,omitempty"`
	ReceiptBASE64 string  `json:"receipt"`
}
