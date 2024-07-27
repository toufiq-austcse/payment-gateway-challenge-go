package tests

import (
	"encoding/json"
	"fmt"
	"github.com/cko-recruitment/payment-gateway-challenge-go/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ginEngine = gin.Default()
var paymentRouterGroup = ginEngine.Group("api/v1/payments")

func TestCreatePayment(t *testing.T) {
	paymentRouterGroup.POST("", handlers.CreatePayment)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/payments", nil)
	if err != nil {
		t.Error("request error", err)
	}
	recorder := httptest.NewRecorder()
	ginEngine.ServeHTTP(recorder, req)

	responseInterface := make(map[string]interface{})
	fmt.Println(recorder.Code)

	unmarshalErr := json.Unmarshal(recorder.Body.Bytes(), &responseInterface)

	if unmarshalErr != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

}
