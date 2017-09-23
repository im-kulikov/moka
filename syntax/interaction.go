package syntax

import "github.com/gcapizzi/moka"

type Interaction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func ReceiveCallTo(methodName string) Interaction {
	return Interaction{methodName: methodName}
}

func (i Interaction) With(args ...interface{}) Interaction {
	return Interaction{methodName: i.methodName, args: args, returnValues: i.returnValues}
}

func (i Interaction) AndReturn(returnValues ...interface{}) Interaction {
	return Interaction{methodName: i.methodName, args: i.args, returnValues: returnValues}
}

func (i Interaction) Apply(double moka.Double) {
	double.StubMethod(i.methodName, i.args, i.returnValues)
}
