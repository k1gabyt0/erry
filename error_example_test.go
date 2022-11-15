package erry_test

import (
	"errors"
	"fmt"

	"github.com/k1gabyt0/erry"
)

func ExampleNewError() {
	errA := errors.New("err A")
	errB := errors.New("err B")

	// Creates brand new multi-error with passed error message
	// and errA and errB ass inner errors.
	multierr := erry.NewError("multierror", errA, errB)
	fmt.Println(multierr)

	if errors.Is(multierr, errA) {
		fmt.Println("This is error A")

	}
	if errors.Is(multierr, errB) {
		fmt.Println("This is error B")

	}
	// Output: multierror:
	//	err A
	//	err B
	// This is error A
	// This is error B
}

func ExampleErrorFrom() {
	errA := errors.New("err A")
	errB := errors.New("err B")

	// Transforms errA into multi-error with errB
	// as one of inner errors.
	multierr := erry.ErrorFrom(errA, errB)
	fmt.Println(multierr)

	if errors.Is(multierr, errA) {
		fmt.Println("This is error A")

	}
	if errors.Is(multierr, errB) {
		fmt.Println("This is error B")

	}
	// Output: err A:
	//	err B
	// This is error A
	// This is error B
}
