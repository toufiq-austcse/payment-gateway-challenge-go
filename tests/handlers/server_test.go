package handlers_test

import (
	"github.com/cko-recruitment/payment-gateway-challenge-go/handlers"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Server", func() {
	var ginEngine = gin.Default()
	var paymentRouterGroup = ginEngine.Group("api/v1/payments")
	var testingServer *httptest.Server

	BeforeEach(func() {
		testingServer = httptest.NewServer(ginEngine)
		//ginEngine.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "api/v1/payments", nil))
		paymentRouterGroup.POST("", handlers.CreatePayment)
	})

	AfterEach(func() {
		testingServer.Close()
	})

	Context("When invalid data is provided", func() {
		It("should return 400 Bad Request", func() {
			response, err := http.Post(testingServer.URL+"/api/v1/payments", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.StatusCode).Should(Equal(400))
		})
	})

})
