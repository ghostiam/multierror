package multierror

import (
	"fmt"
	"io"
	"sync"
)

// Builder contained errors and supports thread-safe addition of errors
type Builder struct {
	mx   sync.Mutex
	errs multiError
}

// Add error (thread-safe), if err is nil, nothing will be added.
func (m *Builder) Add(err error) {
	if err == nil {
		return
	}

	m.mx.Lock()
	m.errs = append(m.errs, err)
	m.mx.Unlock()
}

// ToError returns one error contained many errors or nil if no errors
func (m *Builder) ToError() error {
	if len(m.errs) == 0 {
		return nil
	}

	return m.errs
}

// Errors return the underlying errors, if possible.
// An error value has a errors if it implements the following
// interface:
//
//     type errors interface {
//            Errors() []error
//     }
//
// If the error does not implement Errors, the original error will
// be returned like []error{err}. If the error is nil, nil will be returned without further
// investigation.
func Errors(err error) []error {
	type errors interface {
		Errors() []error
	}

	if err == nil {
		return nil
	}

	errs, ok := err.(errors)
	if !ok {
		return []error{err}
	}

	return errs.Errors()
}

type multiError []error

func (m multiError) Error() string {
	return fmt.Sprintf("%s", m)
}

func (m multiError) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, "MultiError:\n")

	for k, v := range m {
		switch verb {
		case 'v':
			if s.Flag('+') {
				_, _ = fmt.Fprintf(s, "### %d ###\n", k+1)
				_, _ = fmt.Fprintf(s, "%+v\n", v)
				continue
			}

			_, _ = fmt.Fprintf(s, "# %d: ", k+1)
			_, _ = fmt.Fprintf(s, "%v\n", v)
		case 's', 'q':
			_, _ = io.WriteString(s, v.Error())
			_, _ = io.WriteString(s, "\n")
		}
	}
}

func (m multiError) Errors() []error {
	return m
}
