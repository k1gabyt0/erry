package erry

import (
	"errors"
	"fmt"
	"strings"
)

// MError is an error that able to store multiple errors.
type MError struct {
	// original is an error that was passed to be able to store multiple errors.
	original error
	msg      string
	errs     []error
}

// Make sure that MError is an error.
var _ error = &MError{}

// NewError returns MError with passed msg and errs.
//
// Passed nil errs though will be filtered. Only non-nil errs are stored.
func NewError(msg string, errs ...error) *MError {
	e := &MError{
		msg: msg,
	}
	return e.WithErrors(errs...)
}

// ErrorFrom creates an MError err from passed original error
// and it's inner errs.
//
// err's msg will be the same as original.Error().
// Function call errors.Is(err, original) will return true.
//
// If passed original is already an MError, then original is returned.
func ErrorFrom(original error) *MError {
	if original == nil {
		return NewError("")
	}

	if tgt, ok := original.(*MError); ok {
		return tgt
	}

	return &MError{
		original: original,
		msg:      original.Error(),
	}
}

// Error formats error messages with new lines starting with v.msg.
func (v *MError) Error() string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString(v.msg)

	if len(v.errs) > 0 {
		msgBuilder.WriteString(":")

		for _, err := range v.errs {
			msgBuilder.WriteString(fmt.Sprintf("\n\t%s", err.Error()))
		}
	}

	return msgBuilder.String()
}

// Is reports whether any error in v's tree matches target.
//
// It is always returns true if checked against non-nil original.
func (v *MError) Is(target error) bool {
	if errors.Is(v.original, target) {
		return true
	}

	if tgt, ok := target.(*MError); ok {
		isSameMessage := v.msg == tgt.msg
		hasSameErrors := equalErrors(v.errs, tgt.errs)

		if isSameMessage && hasSameErrors {
			return true
		}
	}

	for _, err := range v.errs {
		if is := errors.Is(err, target); is {
			return true
		}
	}
	return false
}

// As finds the first error in v's tree that matches target(starting with the root),
// and if one is found, sets target to that error value and returns true.
// Otherwise, it returns false.
//
// As finds the first matching error in a preorder traversal of the tree.
func (v *MError) As(target any) bool {
	if errors.As(v.original, target) {
		return true
	}

	for _, err := range v.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

// WithErrors stores passed errs in MError.
//
// Passed nil errs though will be filtered. Only non-nil errs are stored.
func (v *MError) WithErrors(errs ...error) *MError {
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}
	v.errs = nonNilErrs
	return v
}

// Message returns root message of MError that will be displayed first.
func (v *MError) Message() string {
	return v.msg
}

// Original error that have been passed to become mutli-error.
func (v *MError) Original() error {
	return v.original
}

// Errors returns all errors stored in v.
func (v *MError) Errors() []error {
	return v.errs
}

func equalErrors(s1, s2 []error) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
