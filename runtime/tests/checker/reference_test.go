package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/language/runtime/sema"
	. "github.com/dapperlabs/flow-go/language/runtime/tests/utils"
)

func TestCheckReferenceTypeSubTyping(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource interface RI {}

          resource R: RI {}

          let ref = &storage[R] as RI
          let ref2: &RI = ref
        `,
	)

	require.Nil(t, err)
}

func TestCheckInvalidReferenceTypeSubTyping(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource interface RI {}

          // NOTE: R does not conform to RI
          resource R {}

          let ref = &storage[R] as RI
          let ref2: &RI = ref
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckReferenceTypeOuter(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource R {}

      fun test(r: &[R]) {}
    `)

	assert.Nil(t, err)
}

func TestCheckReferenceTypeInner(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource R {}

      fun test(r: [&R]) {}
    `)

	assert.Nil(t, err)
}

func TestCheckNestedReferenceType(t *testing.T) {

	_, err := ParseAndCheck(t, `
      resource R {}

      fun test(r: &[&R]) {}
    `)

	assert.Nil(t, err)
}

func TestCheckInvalidReferenceType(t *testing.T) {

	_, err := ParseAndCheck(t, `
      fun test(r: &R) {}
    `)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckReferenceExpressionWithResourceResultType(t *testing.T) {

	checker, err := ParseAndCheckWithStorage(t, `
          resource R {}

          let ref = &storage[R] as R
        `,
	)

	require.Nil(t, err)

	refValueType := checker.GlobalValues["ref"].Type

	assert.IsType(t,
		&sema.ReferenceType{},
		refValueType,
	)

	assert.IsType(t,
		&sema.CompositeType{},
		refValueType.(*sema.ReferenceType).Type,
	)
}

func TestCheckReferenceExpressionWithResourceInterfaceResultType(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource interface T {}
          resource R: T {}

          let ref = &storage[R] as T
        `,
	)

	assert.Nil(t, err)
}

func TestCheckInvalidReferenceExpressionType(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          let ref = &storage[R] as X
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckInvalidReferenceExpressionStorageIndexType(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          let ref = &storage[X] as R
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
}

func TestCheckInvalidReferenceExpressionNonResourceReferencedType(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          struct R {}
          resource T {}

          let ref = &storage[R] as T
        `,
	)

	errs := ExpectCheckerErrors(t, err, 2)

	assert.IsType(t, &sema.NonResourceReferenceError{}, errs[0])
	assert.IsType(t, &sema.TypeMismatchError{}, errs[1])
}

func TestCheckInvalidReferenceExpressionNonResourceResultType(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}
          struct T {}

          let ref = &storage[R] as T
        `,
	)

	errs := ExpectCheckerErrors(t, err, 2)

	assert.IsType(t, &sema.NonResourceReferenceError{}, errs[0])
	assert.IsType(t, &sema.TypeMismatchError{}, errs[1])
}

func TestCheckInvalidReferenceExpressionNonResourceTypes(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          struct R {}
          struct T {}

          let ref = &storage[R] as T
        `,
	)

	errs := ExpectCheckerErrors(t, err, 3)

	assert.IsType(t, &sema.NonResourceReferenceError{}, errs[0])
	assert.IsType(t, &sema.NonResourceReferenceError{}, errs[1])
	assert.IsType(t, &sema.TypeMismatchError{}, errs[2])
}

func TestCheckInvalidReferenceExpressionTypeMismatch(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}
          resource T {}

          let ref = &storage[R] as T
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
}

func TestCheckInvalidReferenceToNonIndex(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          let r <- create R()
          let ref = &r as R
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NonStorageReferenceError{}, errs[0])
}

func TestCheckInvalidReferenceToNonStorage(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          let rs <- [<-create R()]
          let ref = &rs[0] as R
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NonStorageReferenceError{}, errs[0])
}

func TestCheckReferenceUse(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {
              var x: Int

              init() {
                  self.x = 0
              }

              fun setX(_ newX: Int) {
                  self.x = newX
              }
          }

          fun test(): [Int] {
              var r: <-R? <- create R()
              storage[R] <-> r
              // there was no old value, but it must be discarded
              destroy r

              let ref = &storage[R] as R
              ref.x = 1
              let x1 = ref.x
              ref.setX(2)
              let x2 = ref.x
              return [x1, x2]
          }
        `,
	)

	assert.Nil(t, err)
}

func TestCheckReferenceUseArray(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {
              var x: Int

              init() {
                  self.x = 0
              }

              fun setX(_ newX: Int) {
                  self.x = newX
              }
          }

          fun test(): [Int] {
              var rs: <-[R]? <- [<-create R()]
              storage[[R]] <-> rs
              // there was no old value, but it must be discarded
              destroy rs

              let ref = &storage[[R]] as [R]
              ref[0].x = 1
              let x1 = ref[0].x
              ref[0].setX(2)
              let x2 = ref[0].x
              return [x1, x2]
          }
        `,
	)

	assert.Nil(t, err)
}

func TestCheckReferenceIndexingIfReferencedIndexable(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          fun test() {
              var rs: <-[R]? <- [<-create R()]
              storage[[R]] <-> rs
              // there was no old value, but it must be discarded
              destroy rs

              let ref = &storage[[R]] as [R]
              var other <- create R()
              ref[0] <-> other
              destroy other
          }
        `,
	)

	assert.Nil(t, err)
}

func TestCheckInvalidReferenceResourceLoss(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          fun test() {
              var rs: <-[R]? <- [<-create R()]
              storage[[R]] <-> rs
              // there was no old value, but it must be discarded
              destroy rs

              let ref = &storage[[R]] as [R]
              ref[0]
          }
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.ResourceLossError{}, errs[0])
}

func TestCheckInvalidReferenceIndexingIfReferencedNotIndexable(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource R {}

          fun test() {
              var r: <-R? <- create R()
              storage[R] <-> r
              // there was no old value, but it must be discarded
              destroy r

              let ref = &storage[R] as R
              ref[0]
          }
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotIndexableTypeError{}, errs[0])
}

func TestCheckResourceInterfaceReferenceFunctionCall(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource interface I {
              fun foo()
          }

          resource R: I {
              fun foo() {}
          }

          fun test() {
              var r: <-R? <- create R()
              storage[R] <-> r
              // there was no old value, but it must be discarded
              destroy r

              let ref = &storage[R] as I
              ref.foo()
          }
        `,
	)

	assert.Nil(t, err)
}

func TestCheckInvalidResourceInterfaceReferenceFunctionCall(t *testing.T) {

	_, err := ParseAndCheckWithStorage(t, `
          resource interface I {}

          resource R: I {
              fun foo() {}
          }

          fun test() {
              var r: <-R? <- create R()
              storage[R] <-> r
              // there was no old value, but it must be discarded
              destroy r

              let ref = &storage[R] as I
              ref.foo()
          }
        `,
	)

	errs := ExpectCheckerErrors(t, err, 1)

	assert.IsType(t, &sema.NotDeclaredMemberError{}, errs[0])
}