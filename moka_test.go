package moka_test

import (
	"github.com/gcapizzi/moka"
	. "github.com/gcapizzi/moka/syntax"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Moka", func() {
	var collaborator CollaboratorDouble
	var subject Subject

	var failHandlerCalled bool
	var failHandlerMessage string

	BeforeEach(func() {
		failHandlerCalled = false
		failHandlerMessage = ""
		moka.RegisterDoublesFailHandler(func(message string, _ ...int) {
			failHandlerCalled = true
			failHandlerMessage = message
		})

		collaborator = NewCollaboratorDouble()
		subject = NewSubject(collaborator)
	})

	It("supports allowing a method call on a double", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("Query").With("arg").AndReturn("result"))

		Expect(failHandlerCalled).To(BeFalse())

		result := subject.DelegateQuery("arg")

		Expect(result).To(Equal("result"))
	})

	It("makes tests fail on unexpected interactions", func() {
		collaborator.Query("unexpected")

		Expect(failHandlerCalled).To(BeTrue())
		Expect(failHandlerMessage).To(Equal("Unexpected interaction: Query(\"unexpected\")"))
	})

	It("supports expecting a method call on a double", func() {
		ExpectDouble(collaborator).To(ReceiveCallTo("Command").With("arg").AndReturn("result", nil))

		Expect(failHandlerCalled).To(BeFalse())

		result, _ := subject.DelegateCommand("arg")

		Expect(result).To(Equal("result"))
		VerifyCalls(collaborator)
	})
})

type Collaborator interface {
	Query(string) string
	Command(string) (string, error)
}

type CollaboratorDouble struct {
	moka.Double
}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{Double: moka.NewStrictDoubleWithTypeOf(CollaboratorDouble{})}
}

func (d CollaboratorDouble) Query(arg string) string {
	returnValues, err := d.Call("Query", arg)
	if err != nil {
		return ""
	}

	return returnValues[0].(string)
}

func (d CollaboratorDouble) Command(arg string) (string, error) {
	returnValues, err := d.Call("Command", arg)
	if err != nil {
		return "", nil
	}

	returnedString, _ := returnValues[0].(string)
	returnedError, _ := returnValues[1].(error)

	return returnedString, returnedError
}

type Subject struct {
	collaborator Collaborator
}

func NewSubject(collaborator Collaborator) Subject {
	return Subject{collaborator: collaborator}
}

func (s Subject) DelegateQuery(arg string) string {
	return s.collaborator.Query(arg)
}

func (s Subject) DelegateCommand(arg string) (string, error) {
	return s.collaborator.Command(arg)
}
