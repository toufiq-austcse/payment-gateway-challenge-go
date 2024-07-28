package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cko-recruitment/payment-gateway-challenge-go/apimodels/req"
	"github.com/cko-recruitment/payment-gateway-challenge-go/enums"
	"github.com/cko-recruitment/payment-gateway-challenge-go/handlers"
	"github.com/cko-recruitment/payment-gateway-challenge-go/pkg/api_response"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	suite.ginEngine.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/payments", nil))
}

func (suite *integrationTestSuite) TearDownSuite() {
	suite.testingServer.Close()

}
func (suite *integrationTestSuite) Test_CreatePaymentRejectedStatus() {

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

	suite.Run("When all card_number characters is not numeric it should return 400", func() {
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
		body := map[string]interface{}{
			"card_number":      "2222405343248877",
			"expiration_month": 4,
			"expiration_year":  2023,
			"currency":         "USD",
			"amount":           "100",
			"cvv":              "123",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When cvv is missing it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When cvv less than 3 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
			CVV:             "1",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When cvv grater than 4 character it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
			CVV:             "12345",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		suite.Equal(http.StatusBadRequest, response.StatusCode)
	})

	suite.Run("When all cvv characters is not numeric it should return 400", func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
			ExpirationMonth: 4,
			ExpirationYear:  2023,
			Currency:        "USD",
			Amount:          100,
			CVV:             "12a",
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

func (suite *integrationTestSuite) Test_CreatePaymentDeclinedStatus() {
	suite.Run("When Acquiring bank not authorized a payment it should return 200 OK with payment status "+enums.DECLIEND, func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248112",
			ExpirationMonth: 1,
			ExpirationYear:  2026,
			Currency:        "USD",
			Amount:          6000,
			CVV:             "456",
		}
		requestBody, err := json.Marshal(body)
		suite.NoError(err, "no error when marshalling the request")

		response, err := http.Post(suite.testingServer.URL+"/api/v1/payments", "application/json", bytes.NewBuffer(requestBody))
		suite.NoError(err, "no error when calling the endpoint")

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		//var paymentDetails res.PaymentDetails
		var apiBody api_response.Response

		jsonParseError := json.Unmarshal(bodyBytes, &apiBody)
		suite.NoError(jsonParseError, "no error when calling json decode")

		apiResDataMap := apiBody.Data.(map[string]interface{})

		suite.Equal(http.StatusOK, response.StatusCode)
		suite.Equal(enums.DECLIEND, apiResDataMap["status"])

	})

}

func (suite *integrationTestSuite) Test_CreatePaymentAuthorizedStatus() {
	suite.Run("When Acquiring bank authorized a payment it should return 200 OK with payment status "+enums.AUTHORIZED, func() {
		body := req.CreatePaymentReqModel{
			CardNumber:      "2222405343248877",
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

		bodyBytes, _ := io.ReadAll(response.Body)
		print(string(bodyBytes))

		//var paymentDetails res.PaymentDetails
		var apiBody api_response.Response

		jsonParseError := json.Unmarshal(bodyBytes, &apiBody)
		suite.NoError(jsonParseError, "no error when calling json decode")

		apiResDataMap := apiBody.Data.(map[string]interface{})

		suite.Equal(http.StatusOK, response.StatusCode)
		suite.Equal(enums.AUTHORIZED, apiResDataMap["status"])

	})

}

func TestIntegrationTestSuite(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("could not load .env file: %v", err)
	}
	fmt.Println("env loaded")

	suite.Run(t, new(integrationTestSuite))
}
