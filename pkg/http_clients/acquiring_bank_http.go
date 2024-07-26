package http_clients

import (
	"github.com/cko-recruitment/payment-gateway-challenge-go/enums"
	"github.com/go-resty/resty/v2"
	"os"
)

type AcquiringBankResponse struct {
	Authorized        bool   `json:"authorized"`
	AuthorizationCode string `json:"authorization_code"`
}

type AuthorizedPaymentResponse struct {
	ID                  string `json:"id"`
	Status              string `json:"status"`
	LastFourDigitOfCard string `json:"last_four_digit_of_card"`
	ExpiryMonth         int    `json:"expiry_month"`
	ExpiryYear          int    `json:"expiry_year"`
	Currency            string `json:"currency"`
	Amount              int    `json:"amount"`
}

func getBankHttpRequest() *resty.Request {
	return resty.New().SetBaseURL(os.Getenv("ACQUIRING_BANK_BASE_URL")).R().
		SetHeader("Content-Type", "application/json")

}

func AuthorizePayment(cardNumber string, expiryDate string, currency string, amount int, cvv string) (string, error) {
	var apiResponse *AcquiringBankResponse

	resp, err := getBankHttpRequest().
		SetBody(map[string]interface{}{
			"card_number": cardNumber,
			"expiry_date": expiryDate,
			"currency":    currency,
			"amount":      amount,
			"cvv":         cvv,
		}).SetResult(&apiResponse).Post("/payments")

	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {

		return enums.DECLIEND, nil
	}

	if apiResponse.Authorized {
		return enums.AUTHORIZED, nil
	} else {
		return enums.DECLIEND, nil
	}

}
