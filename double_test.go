package moka_test

import (
	"fmt"

	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFailHandlerInvoked bool
var testFailMessage string

func testFailHandler(message string) {
	if !testFailHandlerInvoked {
		testFailHandlerInvoked = true
		testFailMessage = message
	}
}

func resetTestFail() {
	testFailHandlerInvoked = false
	testFailMessage = ""
}

var _ = Describe("StrictDouble", func() {
	var double *StrictDouble
	var returnValues []interface{}

	BeforeEach(func() {
		resetTestFail()
		double = NewStrictDoubleWithFailHandler(testFailHandler)
	})

	Context("when a stubbed method is called", func() {
		BeforeEach(func() {
			double.StubMethod(
				"UltimateQuestion",
				[]interface{}{"life", "universe", "everything"},
				[]interface{}{42, nil},
			)
		})

		Context("with the right arguments", func() {
			BeforeEach(func() {
				returnValues = double.Call("UltimateQuestion", "life", "universe", "everything")
			})

			It("returns the stubbed return values", func() {
				Expect(returnValues).To(Equal([]interface{}{42, nil}))
				Expect(testFailHandlerInvoked).To(BeFalse())
			})
		})

		Context("with the wrong arguments", func() {
			BeforeEach(func() {
				returnValues = double.Call("UltimateQuestion", "foo", "bar")
			})

			It("returns nil and makes the test fail", func() {
				Expect(returnValues).To(BeNil())
				Expect(testFailHandlerInvoked).To(BeTrue())
				Expect(testFailMessage).To(Equal("No stub or mock for method 'UltimateQuestion' with arguments [foo bar]"))
			})
		})
	})

	Context("when a method is mocked", func() {
		BeforeEach(func() {
			double.MockMethod(
				"MakeMeASandwich",
				[]interface{}{"bacon", "lettuce", "tomatoes"},
				[]interface{}{fmt.Errorf("ain't got no bacon mate")},
			)
		})

		Context("and it is called", func() {
			Context("with the right arguments", func() {
				BeforeEach(func() {
					returnValues = double.Call("MakeMeASandwich", "bacon", "lettuce", "tomatoes")
					double.VerifyCalls()
				})

				It("returns the mocked return values and records the call", func() {
					Expect(returnValues).To(Equal([]interface{}{fmt.Errorf("ain't got no bacon mate")}))
					Expect(testFailHandlerInvoked).To(BeFalse())
				})
			})

			Context("with the wrong arguments", func() {
				BeforeEach(func() {
					returnValues = double.Call("MakeMeASandwich", "peanut butter", "jelly")
					double.VerifyCalls()
				})

				It("returns nil and makes the test fail", func() {
					Expect(returnValues).To(BeNil())
					Expect(testFailHandlerInvoked).To(BeTrue())
					Expect(testFailMessage).To(Equal("No stub or mock for method 'MakeMeASandwich' with arguments [peanut butter jelly]"))
				})
			})
		})

		Context("and it is not called", func() {
			BeforeEach(func() {
				double.VerifyCalls()
			})

			It("makes the test fail", func() {
				Expect(testFailHandlerInvoked).To(BeTrue())
				Expect(testFailMessage).To(Equal("Expected the method 'MakeMeASandwich' to be called with arguments [bacon lettuce tomatoes]"))
			})
		})
	})

	Context("when an unknown method is called", func() {
		It("returns nil and makes the test fail", func() {
			returnValues := double.Call("UnstubbedMethod")

			Expect(returnValues).To(BeNil())
			Expect(testFailHandlerInvoked).To(BeTrue())
			Expect(testFailMessage).To(Equal("No stub or mock for method 'UnstubbedMethod' with arguments []"))
		})
	})
})