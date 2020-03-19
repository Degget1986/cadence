package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/cadence/runtime/sema"
	. "github.com/dapperlabs/cadence/runtime/tests/utils"
)

func TestCheckBoolean(t *testing.T) {

	checker, err := ParseAndCheck(t, `
        let x = true
    `)

	require.NoError(t, err)

	assert.Equal(t,
		&sema.BoolType{},
		checker.GlobalValues["x"].Type,
	)
}
