package res

type PaymentDetails struct {
	Id                string `json:"id"`
	Status            string `json:"status"`
	LastFourCardDigit string `json:"last_four_card_digit"`
	ExpiryMonth       int    `json:"expiry_month"`
	ExpiryYear        int    `json:"expiry_year"`
	CurrencyCode      string `json:"currency_code"`
	Amount            int    `json:"amount"`
}

type ProcessPaymentRes struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
