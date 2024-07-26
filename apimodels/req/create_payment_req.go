package req

import (
	"github.com/gin-gonic/gin"
)

type CreatePaymentReqModel struct {
	CardNumber      string `json:"card_number" binding:"required,gte=14,lte=19,number"`
	ExpirationMonth int    `json:"expiration_month" binding:"required,gte=1,lte=12"`
	ExpirationYear  int    `json:"expiration_year" binding:"required"`
	Currency        string `json:"currency" binding:"required,iso4217"`
	Amount          int    `json:"amount" binding:"required"`
	CVV             string `json:"cvv" binding:"required,number,gte=3,lte=4"`
}

func (model *CreatePaymentReqModel) Validate(c *gin.Context) error {
	err := c.ShouldBindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
