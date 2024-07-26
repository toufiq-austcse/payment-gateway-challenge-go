package mapper

import (
	"github.com/cko-recruitment/payment-gateway-challenge-go/apimodels/res"
	"github.com/cko-recruitment/payment-gateway-challenge-go/models"
)

func ToPaymentDetailsRes(payment models.Payment) res.PaymentDetails {
	return res.PaymentDetails{
		Id:                payment.Id,
		Status:            payment.Status,
		LastFourCardDigit: payment.CardNumber[len(payment.CardNumber)-4:],
		ExpiryMonth:       payment.ExpirationMonth,
		ExpiryYear:        payment.ExpirationYear,
		CurrencyCode:      payment.CurrencyCode,
		Amount:            payment.Amount,
	}
}

func ToPaymentModel(id string, status string, cardNumber string, expiryMonth int, expiryYear int, currencyCode string, amount int) models.Payment {
	return models.Payment{
		Id:              id,
		Status:          status,
		CardNumber:      cardNumber,
		ExpirationMonth: expiryMonth,
		ExpirationYear:  expiryYear,
		CurrencyCode:    currencyCode,
		Amount:          amount,
	}
}
