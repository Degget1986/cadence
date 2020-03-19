package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/cadence/runtime/sema"
	. "github.com/dapperlabs/cadence/runtime/tests/utils"
)

func TestCheckCharacterLiteral(t *testing.T) {

	checker, err := ParseAndCheck(t, `
        let a: Character = "a"
    `)

	require.NoError(t, err)

	assert.Equal(t,
		&sema.CharacterType{},
		checker.GlobalValues["a"].Type,
	)
}

func TestCheckInvalidCharacterLiteral(t *testing.T) {

	_, err := ParseAndCheck(t, `
        let a: Character = "abc"
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.InvalidCharacterLiteralError{}, errs[0])
}
