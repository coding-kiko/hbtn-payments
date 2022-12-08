package errors

import (
	"errors"
	"fmt"
)

// Not found error
type NotFound struct {
	Err error
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewNotFound(msg ...string) *NotFound {
	defaultMsg := "payment not found"
	if len(msg) == 0 {
		return &NotFound{Err: errors.New(defaultMsg)}
	}
	return &NotFound{Err: errors.New(msg[0])}
}

// ------------------------------------------------------------

// Not file processing error
type FileError struct {
	Err error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewFileError(msg ...string) *FileError {
	defaultMsg := "error processing file"
	if len(msg) == 0 {
		return &FileError{Err: errors.New(defaultMsg)}
	}
	return &FileError{Err: errors.New(msg[0])}
}

// method not allowed
type MethodNotAllowed struct {
	Err error
}

func (e *MethodNotAllowed) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewMethodNotAllowed(msg ...string) *MethodNotAllowed {
	defaultMsg := "method not allowed"
	if len(msg) == 0 {
		return &MethodNotAllowed{Err: errors.New(defaultMsg)}
	}
	return &MethodNotAllowed{Err: errors.New(msg[0])}
}

// ---

// bad request
type BadRequest struct {
	Err error
}

func (e *BadRequest) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewBadRequest(msg ...string) *BadRequest {
	defaultMsg := "bad request"
	if len(msg) == 0 {
		return &BadRequest{Err: errors.New(defaultMsg)}
	}
	return &BadRequest{Err: errors.New(msg[0])}
}
