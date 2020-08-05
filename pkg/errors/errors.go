package errors

import "fmt"

// https://biosphere.cc/software-engineering/returning-multiple-errors-in-go/
type MultiError struct {
	errs []error
}

func (m MultiError) Errors() []error {
	return m.errs
}

func (m MultiError) Len() int {
	return len(m.errs)
}

// Adds one or many errors
// nil errors are not appended
func (m *MultiError) Add(es ...error) {
	for _, e := range es {
		if e != nil {
			m.errs = append(m.errs, e)
		}
	}
}

// Returns itself if errors are set, otherwise nil.
func (m *MultiError) Return() error {
	if len(m.errs) > 0 {
		return m
	} else {
		return nil
	}
}

func (m MultiError) Error() string {
	switch len(m.errs) {
	case 0:
		return "(0 errors)"
	case 1:
		return m.errs[0].Error()
	case 2:
		return m.errs[0].Error() + " (and 1 other error)"
	default:
		return fmt.Sprintf("%s (and %d other errors)",
			m.errs[0].Error(), len(m.errs)-1)
	}
}
