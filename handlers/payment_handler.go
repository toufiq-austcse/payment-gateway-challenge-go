package handlers

import (
	"errors"
	"github.com/cko-recruitment/payment-gateway-challenge-go/apimodels/req"
	"github.com/cko-recruitment/payment-gateway-challenge-go/enums"
	"github.com/cko-recruitment/payment-gateway-challenge-go/mapper"
	"github.com/cko-recruitment/payment-gateway-challenge-go/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/pkg/api_response"
	"github.com/cko-recruitment/payment-gateway-challenge-go/pkg/http_clients"
	"github.com/cko-recruitment/payment-gateway-challenge-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var paymentMap = make(map[string]models.Payment)

func CreatePayment(context *gin.Context) {
	body := &req.CreatePaymentReqModel{}
	ID, uuidErr := utils.GenerateUUID()
	if uuidErr != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, "Internal Server Error", uuidErr.Error(), nil)
		context.JSON(errRes.Code, errRes)
		return
	}

	paymentStatus := ""
	message := ""

	err := body.Validate(context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, enums.REJECTED, err.Error(), nil)
		context.JSON(errRes.Code, errRes)
		return
	} else {
		expiryDate, expiryDateErr := BuildExpiryDate(body.ExpirationMonth, body.ExpirationYear)
		if expiryDateErr != nil {
			errRes := api_response.BuildErrorResponse(http.StatusBadRequest, enums.REJECTED, expiryDateErr.Error(), nil)
			context.JSON(errRes.Code, errRes)
			return
		}
		status, authError := http_clients.AuthorizePayment(body.CardNumber, expiryDate, body.Currency, body.Amount, body.CVV)
		if authError != nil {
			errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, "Internal Server Error", authError.Error(), nil)
			context.JSON(errRes.Code, errRes)
			return
		} else {
			paymentStatus = status
		}
	}

	paymentModel := mapper.ToPaymentModel(ID, paymentStatus, body.CardNumber, body.ExpirationMonth, body.ExpirationYear, body.Currency, body.Amount)
	paymentMap[paymentModel.Id] = paymentModel

	paymentDetailRes := mapper.ToPaymentDetailsRes(paymentModel)
	res := api_response.BuildResponse(http.StatusOK, message, paymentDetailRes)
	context.JSON(res.Code, res)
	return
}

func GetPaymentById(context *gin.Context) {
	ID := context.Param("id")
	paymentModel, ok := paymentMap[ID]
	if !ok {
		errRes := api_response.BuildErrorResponse(http.StatusNotFound, "Not Found", "", nil)
		context.JSON(errRes.Code, errRes)
		return
	}
	paymentDetailRes := mapper.ToPaymentDetailsRes(paymentModel)
	res := api_response.BuildResponse(http.StatusOK, "", paymentDetailRes)
	context.JSON(res.Code, res)
	return
}

func BuildExpiryDate(expiryMonth int, expiryYear int) (string, error) {
	currentYear, currentMonth, _ := time.Now().Date()
	currentMonthInt := int(currentMonth)

	if expiryYear < currentYear {
		return "", errors.New("expiry date must be in future date")
	}
	if expiryYear == currentYear {
		if expiryMonth < currentMonthInt {
			return "", errors.New("expiry date must be in future date")
		}
	}

	expiryMonthInString := strconv.Itoa(expiryMonth)
	if expiryMonth < 10 {
		expiryMonthInString = "0" + strconv.Itoa(expiryMonth)
	}

	return expiryMonthInString + "/" + strconv.Itoa(expiryYear), nil
}
