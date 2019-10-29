package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/language/runtime/ast"
	"github.com/dapperlabs/flow-go/language/runtime/interpreter"
	"github.com/dapperlabs/flow-go/language/runtime/sema"
)

func TestAssert(t *testing.T) {

	program := &ast.Program{}

	checker, err := sema.NewChecker(
		program,
		ast.StringLocation(""),
		sema.WithPredeclaredValues(BuiltinFunctions.ToValueDeclarations()),
	)
	require.Nil(t, err)

	inter, err := interpreter.NewInterpreter(
		checker,
		interpreter.WithPredefinedValues(BuiltinFunctions.ToValues()),
	)
	require.Nil(t, err)

	_, err = inter.Invoke("assert", false, "oops")
	assert.Equal(t, err, AssertionError{
		Message:  "oops",
		Location: interpreter.LocationPosition{},
	})

	_, err = inter.Invoke("assert", false)
	assert.Equal(t, err, AssertionError{
		Message:  "",
		Location: interpreter.LocationPosition{},
	})

	_, err = inter.Invoke("assert", true, "oops")
	assert.Nil(t, err)

	_, err = inter.Invoke("assert", true)
	assert.Nil(t, err)
}

func TestPanic(t *testing.T) {

	checker, err := sema.NewChecker(
		&ast.Program{},
		ast.StringLocation(""),
		sema.WithPredeclaredValues(BuiltinFunctions.ToValueDeclarations()),
	)
	require.Nil(t, err)

	inter, err := interpreter.NewInterpreter(
		checker,
		interpreter.WithPredefinedValues(BuiltinFunctions.ToValues()),
	)
	require.Nil(t, err)

	_, err = inter.Invoke("panic", "oops")
	assert.Equal(t, err, PanicError{
		Message:  "oops",
		Location: interpreter.LocationPosition{},
	})
}