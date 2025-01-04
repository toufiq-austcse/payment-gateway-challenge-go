package handlers_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Payment", func() {
	BeforeEach(func() {
		fmt.Println("Before Each Called")
	})

	Context("Create Payment", func() {
		It("200 response", func() {
			Expect(2).Should(Equal(2))
		})
		It("2003 response", func() {
			Expect(2).Should(Equal(2))
		})
	})

})
