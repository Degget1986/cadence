package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/language/runtime/sema"
	. "github.com/dapperlabs/flow-go/language/runtime/tests/utils"
)

func TestCheckCastingIntLiteralToInt8(t *testing.T) {

	checker, err := ParseAndCheck(t, `
      let x = 1 as Int8
    `)

	require.Nil(t, err)

	assert.Equal(t,
		&sema.Int8Type{},
		checker.GlobalValues["x"].Type,
	)

	assert.NotEmpty(t, checker.Elaboration.CastingTargetTypes)
}

func TestCheckInvalidCastingIntLiteralToString(t *testing.T) {

	_, err := ParseAndCheck(t, `
      let x = 1 as String
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckCastingIntLiteralToAny(t *testing.T) {

	checker, err := ParseAndCheck(t, `
      let x = 1 as Any
    `)

	require.Nil(t, err)

	assert.Equal(t,
		&sema.AnyType{},
		checker.GlobalValues["x"].Type,
	)

	assert.NotEmpty(t, checker.Elaboration.CastingTargetTypes)
}

func TestCheckCastingArrayLiteral(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun zipOf3(a: [Any; 3], b: [Int; 3]): [[Any; 2]; 3] {
          return [
              [a[0], b[0]] as [Any; 2],
              [a[1], b[1]] as [Any; 2],
              [a[2], b[2]] as [Any; 2]
          ]
      }
    `)

	require.Nil(t, err)
}