package bot

type GetPaymentSummaryResponse struct {
	SummaryBase64 string `json:"summary"`
}

type ReqBody struct {
	Month         string  `json:"month"`
	Amount        int     `json:"amount"`
	EmailTo       *string `json:"emailto"`
	Company       string  `json:"company"`
	ReceiptBASE64 string  `json:"receipt"`
}
