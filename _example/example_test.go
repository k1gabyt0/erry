package _example_test

import (
	"errors"
	"fmt"
	"github.com/k1gabyt0/erry"
	"testing"
)

func ExampleNewError() {
	errA := errors.New("err A")
	errB := errors.New("err B")

	// Creates brand new multi-error with passed error message
	// and errA and errB ass inner errors.
	multierr := erry.NewError("multierror", errA, errB)
	fmt.Println(multierr)
	// Output: multierror:
	//           err A
	//           err B
	if errors.Is(multierr, errA) {
		fmt.Println("This is error A")
		// Output: This is error A
	}
	if errors.Is(multierr, errB) {
		fmt.Println("This is error B")
		// Output: This is error B
	}
}

func ExampleErrorFrom() {
	errA := errors.New("err A")
	errB := errors.New("err B")

	// Transforms errA into multi-error with errB
	// as one of inner errors.
	multierr := erry.ErrorFrom(errA, errB)
	fmt.Println(multierr)
	// Output: err A:
	//		     err B
	if errors.Is(multierr, errA) {
		fmt.Println("This is error A")
		// Output: This is error A
	}
	if errors.Is(multierr, errB) {
		fmt.Println("This is error B")
		// Output: This is error B
	}
}
