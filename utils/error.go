package utils

import "errors"

var (
	// ErrIDNotValid is a not valid id error
	ErrIDNotValid = errors.New("ID must be valid")
	// ErrAlreadyPaid is an already paid for parking error
	ErrAlreadyPaid = errors.New("You have already paid")
	// ErrPayFirst needed to pay before checking out
	ErrPayFirst = errors.New("You have to pay first")
	// ErrAlreadyCheckedOut is an already checked out error
	ErrAlreadyCheckedOut = errors.New("You have already checked out")
	// ErrPlateNotValid is a validation error
	ErrPlateNotValid = errors.New("Plate must be valid, format: AAA-1234")
	// ErrBadRequest is a request-reading error
	ErrBadRequest = errors.New("Bad request")
	// ErrInternalServer is an internal problem
	ErrInternalServer = errors.New("Internal server error")
	// ErrImageRecognition is used when we can't parse image words
	ErrImageRecognition = errors.New("Image recognition failed")
	// ErrNotFound is used when we can't find
	ErrNotFound = errors.New("Not found")
	// ErrMethodNotAllowed is used when a method is not allowed for an endpoint
	ErrMethodNotAllowed = errors.New("Method not allowed")
)
