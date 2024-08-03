package http

import "net/http"

type unAuthorized interface {
	// This error is returned when anything related to authentication fails.
	UnAuthorized()
}

type notFound interface {
	// This error is returned when the resource is not found.
	NotFound()
}

type badRequest interface {
	// This error is returned when the request is invalid. Normally wrong arguments.
	BadRequest()
}

type internalServer interface {
	// This error is returned when the server fails to process the request.
	// It's normally because of errors that are not supposed to happen.
	InternalServer()
}

func getHttpStatus(err error) int {
	// Check auth first. and then argument validation. and then not found. and then internal server problem.
	if isUnAuthorized(err) {
		return http.StatusUnauthorized
	}
	if isBadRequest(err) {
		return http.StatusBadRequest
	}
	if isNotFound(err) {
		return http.StatusNotFound
	}
	if isInternalServerProblem(err) {
		return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}

func isUnAuthorized(err error) bool {
	unAuthorized, ok := err.(unAuthorized)
	return ok && unAuthorized != nil
}

func isNotFound(err error) bool {
	notFound, ok := err.(notFound)
	return ok && notFound != nil
}

func isBadRequest(err error) bool {
	badRequest, ok := err.(badRequest)
	return ok && badRequest != nil
}

func isInternalServerProblem(err error) bool {
	internalServer, ok := err.(internalServer)
	return ok && internalServer != nil
}
