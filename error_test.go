package erry_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/k1gabyt0/erry"
)

func TestMError_Equality(t *testing.T) {
	err1 := erry.NewError("error epta")
	err2 := erry.NewError("error epta")
	if err1 == err2 {
		t.Errorf("errors(%q, %q) with same message should not be equal", err1, err2)
	}
}

func TestMError_New(t *testing.T) {
	type args struct {
		msg  string
		errs []error
	}

	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{
			name: "Message passed",
			args: args{
				msg: "this is an error",
			},
			wantErrs: []error{},
		},
		{
			name: "Empty message is fine too",
			args: args{
				msg: "",
			},
			wantErrs: []error{},
		},
		{
			name: "No errs passed - no errs stored",
			args: args{
				msg:  "this is an error",
				errs: []error{},
			},
			wantErrs: []error{},
		},
		{
			name: "All nil errs passed - no errs stored",
			args: args{
				msg:  "this is an error",
				errs: []error{nil, nil, nil},
			},
			wantErrs: []error{},
		},
		{
			name: "Non nil errs passed - non nil errs stored",
			args: args{
				msg: "this is an error",
				errs: []error{
					err1,
					err2,
					err3,
				},
			},
			wantErrs: []error{
				err1,
				err2,
				err3,
			},
		},
		{
			name: "Non nil & nil errs passed - non nil errs stored",
			args: args{
				msg: "this is an error",
				errs: []error{
					err1,
					nil,
					err2,
					nil,
					err3,
					nil,
				},
			},
			wantErrs: []error{
				err1,
				err2,
				err3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := erry.NewError(tt.args.msg, tt.args.errs...)
			if err == nil {
				t.Error("no MError was created")
				return
			}

			// check message
			if err.Message() != tt.args.msg {
				t.Errorf("expected message=%qm but got=%q", tt.args.msg, err.Message())
			}
			if !strings.Contains(err.Error(), err.Message()) {
				t.Errorf("Error()=%q doesn't contain root's messsage=%q", err.Error(), err.Message())
			}

			// check errs
			errs := err.Errors()
			if len(errs) == len(tt.wantErrs) {
				for _, e := range errs {
					var found bool
					for _, wantE := range tt.wantErrs {
						if wantE == e {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("returned stored error=%q is not int wantErr=%v", e, tt.wantErrs)
					}
				}
			} else {
				t.Errorf("errs we want=%v and errs we got=%v have different sizes", tt.wantErrs, errs)
			}
		})
	}
}

func TestMError_From(t *testing.T) {
	type args struct {
		original error
		errs     []error
	}

	ErrA := errors.New("error A")
	ErrB := errors.New("error B")
	ErrC := errors.New("error C")
	MErrWithA := erry.NewError("validation error", ErrA)

	tests := []struct {
		name        string
		args        args
		wantMsg     string
		wantFullMsg string
		wantIs      []error
		wantIsNot   []error
	}{
		{
			name: "If original is nil, then return empty MError",
			args: args{
				original: nil,
			},
			wantMsg:     "",
			wantFullMsg: "",
			wantIsNot:   []error{ErrA, ErrB, ErrC, MErrWithA},
		},
		{
			name: "Passed error is simple",
			args: args{
				original: ErrA,
			},
			wantMsg:     ErrA.Error(),
			wantFullMsg: ErrA.Error(),
			wantIs:      []error{ErrA},
			wantIsNot:   []error{ErrB, ErrC, MErrWithA},
		},
		{
			name: "Passed error and children",
			args: args{
				original: ErrA,
				errs:     []error{ErrB, ErrC},
			},
			wantMsg:     ErrA.Error(),
			wantFullMsg: fmt.Sprintf("%s:\n\t%s\n\t%s", ErrA.Error(), ErrB.Error(), ErrC),
			wantIs:      []error{ErrA, ErrB, ErrC},
			wantIsNot:   []error{MErrWithA},
		},
		{
			name: "Passed error and children(some of them nil)",
			args: args{
				original: ErrA,
				errs: []error{
					nil,
					ErrB,
					nil,
					ErrC,
					nil,
					nil,
				},
			},
			wantMsg:     ErrA.Error(),
			wantFullMsg: fmt.Sprintf("%s:\n\t%s\n\t%s", ErrA.Error(), ErrB.Error(), ErrC),
			wantIs:      []error{ErrA, ErrB, ErrC},
			wantIsNot:   []error{MErrWithA},
		},
		{
			name: "Passed error is another MError",
			args: args{
				original: MErrWithA,
			},
			wantMsg:     MErrWithA.Message(),
			wantFullMsg: MErrWithA.Error(),
			wantIs:      []error{ErrA, MErrWithA},
			wantIsNot:   []error{ErrB, ErrC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := erry.ErrorFrom(tt.args.original, tt.args.errs...)
			if err == nil {
				t.Error("created MError is nil")
				return
			}

			// check message
			if tt.wantMsg != err.Message() {
				t.Errorf("wanted error message to be=%q, but got=%q", tt.wantMsg, err.Message())
			}
			if tt.wantFullMsg != err.Error() {
				t.Errorf("wanted full error message to be=%q, but got=%q", tt.wantFullMsg, err.Error())
			}

			// check errors
			for _, wantErr := range tt.wantIs {
				if !errors.Is(err, wantErr) {
					t.Errorf("expected that %q is %q", err, wantErr)
				}
			}
			for _, dontWantErr := range tt.wantIsNot {
				if errors.Is(err, dontWantErr) {
					t.Errorf("expected that %q is NOT %q", err, dontWantErr)
				}
			}
		})
	}
}

func TestMError_Is(t *testing.T) {
	errA := errors.New("validation A failed")
	errB := errors.New("validation B failed")
	errC := errors.New("validation C failed")

	errComplexBAndC := erry.NewError("errComplexBAndC", errB, errC)
	errSuperComplex := erry.NewError("errSuperComplex", errComplexBAndC)

	type args struct {
		innerErrs []error
	}

	tests := []struct {
		name          string
		args          args
		wantErrIs     []error
		dontWantErrIs []error
	}{
		{
			name: "No inner errors",
			args: args{
				innerErrs: []error{},
			},
			wantErrIs:     []error{},
			dontWantErrIs: []error{errA, errB, errC},
		},
		{
			name: "One inner error",
			args: args{
				innerErrs: []error{errA},
			},
			wantErrIs:     []error{errA},
			dontWantErrIs: []error{errB, errC},
		},
		{
			name: "Many inner errors",
			args: args{
				innerErrs: []error{errA, errB, errC},
			},
			wantErrIs:     []error{errA, errB, errC},
			dontWantErrIs: []error{},
		},
		{
			name: "Complex inner error with some errors",
			args: args{
				innerErrs: []error{errSuperComplex},
			},
			wantErrIs:     []error{errB, errC, errComplexBAndC, errSuperComplex},
			dontWantErrIs: []error{errA},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := erry.NewError("validation error", tt.args.innerErrs...)
			for _, e := range tt.wantErrIs {
				if !errors.Is(err, e) {
					t.Errorf("expected err=(%q) to be (%q)", err, e)
				}
			}
			for _, e := range tt.dontWantErrIs {
				if errors.Is(err, e) {
					t.Errorf("not expected err=(%q) to be (%q)", err, e)
				}
			}
		})
	}
}

type simpleValidationError struct {
	message string
}

func (v simpleValidationError) Error() string {
	return v.message
}

func TestMError_As(t *testing.T) {
	errB := &simpleValidationError{message: "validation B failed"}
	errC := &simpleValidationError{message: "validation C failed"}
	errComplexBAndC := erry.NewError("errComplexBAndC", errB, errC)

	type args struct {
		original  error
		innerErrs []error
	}

	tests := []struct {
		name          string
		args          args
		targetsOkFn   func() *simpleValidationError
		targetsFailFn func() *simpleValidationError
	}{
		{
			name: "Should not be setted to unrelated type",
			args: args{},
			targetsFailFn: func() *simpleValidationError {
				var unrelated *simpleValidationError
				return unrelated
			},
		},
		{
			name: "Should be setted to related type",
			args: args{
				innerErrs: []error{errB},
			},
			targetsOkFn: func() *simpleValidationError {
				var related *simpleValidationError
				return related
			},
		},
		{
			name: "Should be setted to related type in complex case",
			args: args{
				innerErrs: []error{errComplexBAndC},
			},
			targetsOkFn: func() *simpleValidationError {
				var related *simpleValidationError
				return related
			},
		},
		{
			name: "Should be setted to original type",
			args: args{
				original: errB,
			},
			targetsOkFn: func() *simpleValidationError {
				var original *simpleValidationError
				return original
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.args.original != nil {
				err = erry.ErrorFrom(tt.args.original, tt.args.innerErrs...)
			} else {
				err = erry.NewError("validation error", tt.args.innerErrs...)
			}

			errType := reflect.TypeOf(err)

			// should be setted to itself
			var selfErr *erry.MError
			if !errors.As(err, &selfErr) {
				t.Errorf("expected %q to be setted to itself", errType)
			}

			if tt.targetsOkFn != nil {
				targetOk := tt.targetsOkFn()
				trgtType := reflect.TypeOf(targetOk)
				if !errors.As(err, &targetOk) {
					t.Errorf("expected %q to be setted into %q", errType, trgtType)
				}
			}

			if tt.targetsFailFn != nil {
				targetFail := tt.targetsFailFn()
				trgtType := reflect.TypeOf(targetFail)
				if errors.As(err, &targetFail) {
					t.Errorf("expected %q to be NOT setted into %q", errType, trgtType)
				}
			}
		})
	}
}
