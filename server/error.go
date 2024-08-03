package server

import "fmt"

type UnauthorizedError struct {
	originalError error
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.originalError)
}

func (e *UnauthorizedError) Unwrap() error {
	return e.originalError
}

func (e *UnauthorizedError) UnAuthorized() {}

func newUnauthorizedError(err error) *UnauthorizedError {
	return &UnauthorizedError{originalError: err}
}

type NotFoundError struct {
	originalError error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %v", e.originalError)
}

func (e *NotFoundError) Unwrap() error {
	return e.originalError
}

func (e *NotFoundError) NotFound() {}

func newNotFoundError(err error) *NotFoundError {
	return &NotFoundError{originalError: err}
}

type BadRequestError struct {
	originalError error
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %v", e.originalError)
}

func (e *BadRequestError) Unwrap() error {
	return e.originalError
}

func (e *BadRequestError) BadRequest() {}

func newBadRequestError(err error) *BadRequestError {
	return &BadRequestError{originalError: err}
}

type InternalServerError struct {
	originalError error
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.originalError)
}

func (e *InternalServerError) Unwrap() error {
	return e.originalError
}

func (e *InternalServerError) InternalServer() {}

func newInternalServerError(err error) *InternalServerError {
	return &InternalServerError{originalError: err}
}
