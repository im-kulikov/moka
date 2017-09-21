package moka_test

import (
	"testing"

	"github.com/gcapizzi/moka"
)

func TestQueryDelegation(t *testing.T) {
	m := moka.New(t)

	collaborator := NewCollaboratorDouble()
	m.Allow(collaborator).ToReceive("Query").With("arg").AndReturn("result")

	subject := NewSubject(collaborator)

	queryResult := subject.DelegateQuery("arg")

	if queryResult != "result" {
		t.Errorf("Query result: '%s'", queryResult)
	}
}

type Collaborator interface {
	Query(string) string
}

type CollaboratorDouble struct {
	moka.Double
}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{Double: moka.NewConcreteDouble()}
}

func (d CollaboratorDouble) Query(arg string) string {
	return d.Call("Query", arg).Get(0).(string)
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
