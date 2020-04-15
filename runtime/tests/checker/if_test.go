package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/cadence/runtime/sema"
	. "github.com/dapperlabs/cadence/runtime/tests/utils"
)

func TestCheckIfStatementTest(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          if true {}
      }
    `)

	require.NoError(t, err)
}

func TestCheckIfStatementScoping(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          if true {
              let x = 1
          }
          x
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckInvalidIfStatementTest(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          if 1 {}
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckInvalidIfStatementElse(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          if true {} else {
              x
          }
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckIfStatementTestWithDeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(x: Int?): Int {
          if var y = x {
              return y
          }

          return 0
      }
    `)

	require.NoError(t, err)
}

func TestCheckInvalidIfStatementTestWithDeclarationReferenceInElse(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(x: Int?) {
          if var y = x {
              // ...
          } else {
              y
          }
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckIfStatementTestWithDeclarationNestedOptionals(t *testing.T) {

	_, err := ParseAndCheck(t, `
     fun test(x: Int??): Int? {
         if var y = x {
             return y
         }

         return nil
     }
    `)

	require.NoError(t, err)
}

func TestCheckIfStatementTestWithDeclarationNestedOptionalsExplicitAnnotation(t *testing.T) {

	_, err := ParseAndCheck(t, `
     fun test(x: Int??): Int? {
         if var y: Int? = x {
             return y
         }

         return nil
     }
    `)

	require.NoError(t, err)
}

func TestCheckInvalidIfStatementTestWithDeclarationNonOptional(t *testing.T) {

	_, err := ParseAndCheck(t, `
     fun test(x: Int) {
         if var y = x {
             // ...
         }

         return
     }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckInvalidIfStatementTestWithDeclarationSameType(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(x: Int?): Int? {
          if var y: Int? = x {
             return y
          }

          return nil
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}
