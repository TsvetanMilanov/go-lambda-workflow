package workflow

import (
	"fmt"
	"runtime/debug"
)

// Error describes Workflow error.
type Error interface {
	Error() string
	Stack() string
	OriginalError() error
}

type workflowError struct {
	originalError error
	stack         string
	message       string
}

func (e *workflowError) Error() string {
	return e.message
}

func (e *workflowError) Stack() string {
	return e.stack
}

func (e *workflowError) OriginalError() error {
	return e.originalError
}

func newErrorWithMessage(format string, args ...interface{}) Error {
	if len(args) > 0 {
		// If the args len is 0 and we pass it to the Erorrf func,
		// the error message will contain some more data.
		return newError(fmt.Errorf(format, args...))
	}

	return newError(fmt.Errorf(format))
}

func newError(err error) Error {
	if err == nil {
		return nil
	}

	return &workflowError{originalError: err, message: err.Error(), stack: getStack()}
}

func getStack() string {
	return string(debug.Stack())
}
