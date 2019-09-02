package multierror_test

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/ghostiam/multierror"
)

func TestBuilder_ToError_Empty(t *testing.T) {
	var errs multierror.Builder
	err := errs.ToError()
	if err != nil {
		t.FailNow()
	}
}

func TestBuilder_Add_Concurrency(t *testing.T) {
	var wg sync.WaitGroup
	var errs multierror.Builder

	count := 1000
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			errs.Add(errors.New("test-" + strconv.Itoa(i)))
		}(i)
	}
	wg.Wait()
}

func TestErrors(t *testing.T) {
	err := buildMultiError()
	errs := multierror.Errors(err)
	if len(errs) != 3 {
		t.FailNow()
	}

	expectedErrs := []error{
		errors.New("test1"),
		errors.New("test2"),
		errors.New("test3"),
	}
	if !reflect.DeepEqual(errs, expectedErrs) {
		t.FailNow()
	}
}

func TestErrors_Invalid(t *testing.T) {
	err := errors.New("not multiError")
	errs := multierror.Errors(err)
	if errs != nil && errs[0] != err {
		t.FailNow()
	}
}

func TestErrors_Nil(t *testing.T) {
	errs := multierror.Errors(nil)
	if errs != nil {
		t.FailNow()
	}
}

func TestMultiError_Error(t *testing.T) {
	err := buildMultiError()
	expectedErr := "MultiError:\ntest1\ntest2\ntest3\n"
	if err == nil || err.Error() != expectedErr {
		t.FailNow()
	}
}

func TestMultiError_Format_S(t *testing.T) {
	err := buildMultiError()
	if fmt.Sprintf("%s", err) != "MultiError:\ntest1\ntest2\ntest3\n" {
		t.FailNow()
	}
}

func TestMultiError_Format_V(t *testing.T) {
	err := buildMultiError()
	if fmt.Sprintf("%v", err) != "MultiError:\n# 1: test1\n# 2: test2\n# 3: test3\n" {
		t.FailNow()
	}
}

func TestMultiError_Format_VPlus(t *testing.T) {
	err := buildMultiError()
	if fmt.Sprintf("%+v", err) != "MultiError:\n### 1 ###\ntest1\n### 2 ###\ntest2\n### 3 ###\ntest3\n" {
		t.FailNow()
	}
}

func buildMultiError() error {
	var errs multierror.Builder

	errs.Add(errors.New("test1"))
	errs.Add(errors.New("test2"))
	errs.Add(errors.New("test3"))
	errs.Add(nil)

	return errs.ToError()
}

func BenchmarkBuilder_Add(b *testing.B) {
	var errs multierror.Builder
	err := errors.New("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errs.Add(err)
	}
}
