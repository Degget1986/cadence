package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/cadence/runtime/sema"
	. "github.com/dapperlabs/cadence/runtime/tests/utils"
)

func TestCheckReferenceInFunction(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          test
      }
    `)

	require.NoError(t, err)
}

func TestCheckParameterNameWithFunctionName(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(test: Int) {
          test
      }
    `)

	require.NoError(t, err)
}

func TestCheckMutuallyRecursiveFunctions(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun isEven(_ n: Int): Bool {
          if n == 0 {
              return true
          }
          return isOdd(n - 1)
      }

      fun isOdd(_ n: Int): Bool {
          if n == 0 {
              return false
          }
          return isEven(n - 1)
      }
    `)

	require.NoError(t, err)
}

func TestCheckMutuallyRecursiveScoping(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun f(): Int {
         return g()
      }

      let x = f()
      let y = 0

      fun g(): Int {
          return y
      }
    `)

	require.NoError(t, err)
}

func TestCheckInvalidFunctionDeclarations(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test() {
          fun foo() {}
          fun foo() {}
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.RedeclarationError{}, errs[0])
}

func TestCheckInvalidFunctionRedeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun foo() {
          fun foo() {}
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.RedeclarationError{}, errs[0])
}

func TestCheckFunctionAccess(t *testing.T) {

	_, err := ParseAndCheck(t, `
       pub fun test() {}
    `)

	require.NoError(t, err)
}

func TestCheckInvalidFunctionAccess(t *testing.T) {

	_, err := ParseAndCheck(t, `
       pub(set) fun test() {}
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.InvalidAccessModifierError{}, errs[0])
}

func TestCheckReturnWithoutExpression(t *testing.T) {

	_, err := ParseAndCheck(t, `
       fun returnNothing() {
           return
       }
    `)

	require.NoError(t, err)
}

func TestCheckAnyReturnType(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun foo(): AnyStruct {
          return foo
      }
    `)

	require.NoError(t, err)
}

func TestCheckInvalidParameterTypes(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(x: X, y: Y) {}
    `)

	errs := ExpectCheckerErrors(t, err, 2)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])

	assert.IsType(t, &sema.NotDeclaredError{}, errs[1])

}

func TestCheckInvalidParameterNameRedeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(a: Int, a: Int) {}
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.RedeclarationError{}, errs[0])
}

func TestCheckParameterRedeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(a: Int) {
          let a = 1
      }
    `)

	require.NoError(t, err)
}

func TestCheckInvalidParameterAssignment(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(a: Int) {
          a = 1
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.AssignmentToConstantError{}, errs[0])
}

func TestCheckInvalidArgumentLabelRedeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(x a: Int, x b: Int) {}
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.RedeclarationError{}, errs[0])
}

func TestCheckArgumentLabelRedeclaration(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(_ a: Int, _ b: Int) {}
    `)

	require.NoError(t, err)
}

func TestCheckInvalidFunctionDeclarationReturnValue(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(): Int {
          return true
      }
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckInvalidResourceCapturingThroughVariable(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource Kitty {}

      fun makeKittyCloner(): ((): @Kitty) {
          let kitty <- create Kitty()
          return fun (): @Kitty {
              return <-kitty
          }
      }

      let test = makeKittyCloner()
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.ResourceCapturingError{}, errs[0])
}

func TestCheckInvalidResourceCapturingThroughParameter(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource Kitty {}

      fun makeKittyCloner(kitty: @Kitty): ((): @Kitty) {
          return fun (): @Kitty {
              return <-kitty
          }
      }

      let test = makeKittyCloner(kitty: <-create Kitty())
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.ResourceCapturingError{}, errs[0])
}

func TestCheckInvalidSelfResourceCapturing(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource Kitty {
          fun makeCloner(): ((): @Kitty) {
              return fun (): @Kitty {
                  return <-self
              }
          }
      }

      let kitty <- create Kitty()
      let test = kitty.makeCloner()
    `)

	errs := ExpectCheckerErrors(t, err, 2)

	assert.IsType(t, &sema.ResourceCapturingError{}, errs[0])
	assert.IsType(t, &sema.InvalidSelfInvalidationError{}, errs[1])
}

func TestCheckInvalidResourceCapturingJustMemberAccess(t *testing.T) {
	// Resource capturing even just for read access (e.g. reading a member) is invalid

	_, err := ParseAndCheck(t, `
      resource Kitty {
          let id: Int

          init(id: Int) {
              self.id = id
          }
      }

      fun makeKittyIdGetter(): ((): Int) {
          let kitty <- create Kitty(id: 1)
          let getId = fun (): Int {
              return kitty.id
          }
          destroy kitty
          return getId
      }

      let test = makeKittyIdGetter()
    `)

	errs := ExpectCheckerErrors(t, err, 2)

	assert.IsType(t, &sema.ResourceCapturingError{}, errs[0])
	assert.IsType(t, &sema.ResourceLossError{}, errs[1])
}
