package errors

import "fmt"

// MultiError represents multiple errors
// https://biosphere.cc/software-engineering/returning-multiple-errors-in-go/
type MultiError struct {
	errs []error
}

// Errors returns the list of errors
func (m MultiError) Errors() []error {
	return m.errs
}

// Len returns the number of errors
func (m MultiError) Len() int {
	return len(m.errs)
}

// Add adds one or many errors
// nil errors are not appended
func (m *MultiError) Add(es ...error) {
	for _, e := range es {
		if e != nil {
			m.errs = append(m.errs, e)
		}
	}
}

// Return returns itself if errors are set, otherwise nil.
func (m *MultiError) Return() error {
	if len(m.errs) > 0 {
		return m
	}

	return nil
}

const (
	noErrors       = 0
	singleError    = 1
	exactly2Errors = 2
)

// Error returns a human readable string indicating the number of errors that are contained
// but no specifics on the actual errors
// TODO: make it work like error wrapping
func (m MultiError) Error() string {
	switch len(m.errs) {
	case noErrors:
		return "(0 errors)"
	case singleError:
		return m.errs[0].Error()
	case exactly2Errors:
		return m.errs[0].Error() + " (and 1 other error)"
	default:
		return fmt.Sprintf("%s (and %d other errors)",
			m.errs[0].Error(), len(m.errs)-1)
	}
}
