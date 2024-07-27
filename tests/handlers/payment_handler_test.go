package tests

import (
	"bytes"
	"encoding/json"
	"github.com/cko-recruitment/payment-gateway-challenge-go/apimodels/req"
	"github.com/cko-recruitment/payment-gateway-challenge-go/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type integrationTestSuite struct {
	suite.Suite
	ginEngine          *gin.Engine
	paymentRouterGroup *gin.RouterGroup
	baseUrl            string
	testingServer      *httptest.Server
}

func (suite *integrationTestSuite) SetupSuite() {
	suite.ginEngine = gin.Default()
	suite.paymentRouterGroup = suite.ginEngine.Group("api/v1/payments")
	suite.paymentRouterGroup.POST("", handlers.CreatePayment)
	suite.baseUrl = "http://localhost:8081"

	suite.testingServer = httptest.NewServer(suite.ginEngine)
}

func (suite *integrationTestSuite) TearDownSuite() {
	suite.testingServer.Close()

}
func (suite *integrationTestSuite) Test_CreatePaymentRejectedStatus() {
	suite.ginEngine.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/payments", nil))

	suite.Run("When body is empty it should return 400", func() {
		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer([]byte{}))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When card_number is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			ExpirationMonth: 4,
			ExpirationYear:  2025,
			Currency:        "GBP",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})
	suite.Run("When card_number length is less than 14 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "12345567",
			ExpirationMonth: 4,
			ExpirationYear:  2025,
			Currency:        "GBP",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})
	suite.Run("When card_number length is grater than 19 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "12345567912345678912",
			ExpirationMonth: 4,
			ExpirationYear:  2025,
			Currency:        "GBP",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When card_number characters is not numeric it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "123w45567",
			ExpirationMonth: 4,
			ExpirationYear:  2025,
			Currency:        "GBP",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When expiration_month is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:     "123w45567",
			ExpirationYear: 2025,
			Currency:       "GBP",
			Amount:         100,
			CVV:            "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})
	suite.Run("When expiration_month is less than 1 it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 0,
			ExpirationYear:  2026,
			Currency:        "USD",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When expiration_month is less grater than 12 it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 13,
			ExpirationYear:  2026,
			Currency:        "USD",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When expiration_year is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			Currency:        "USD",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When expiration_year is past it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When currency is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})
	suite.Run("When currency is grater than 3 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USDC",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When currency is less than 3 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "US",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When amount is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When amount is not integer it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
			CVV:             "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(integrationTestSuite))
}
