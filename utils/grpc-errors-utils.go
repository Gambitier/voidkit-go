package utils

import (
	"fmt"
	"strings"

	"github.com/Gambitier/voidkitgo/proto/common"
)

// ResultWithError represents a proto-generated response type that has Success and Error fields
type ResultWithError interface {
	GetSuccess() bool
	GetError() *common.Error
}

// ResponseWithResults represents a proto-generated response type that has a Results field
type ResponseWithResults[T ResultWithError] interface {
	GetResults() []T
}

// ErrorResult stores information about an error in a response
type ErrorResult[T ResultWithError] struct {
	Index    int
	Response T
	Err      *common.Error
}

func (e *ErrorResult[T]) Error() string {
	return fmt.Sprintf("Error in response at index %d: %s", e.Index, e.Err.Message)
}

type ErrorResultList[T ResultWithError] []ErrorResult[T]

func (e *ErrorResultList[T]) Error() string {
	var errorMessages []string
	for _, err := range *e {
		errorMessages = append(errorMessages, err.Error())
	}
	return strings.Join(errorMessages, "\n")
}

// CheckIfAnyErrors checks if any response in an array has an error and returns the error details
func CheckIfAnyErrors[T ResultWithError](responses []T) (bool, ErrorResultList[T]) {
	var errorResults ErrorResultList[T]

	for i, response := range responses {
		if !response.GetSuccess() {
			errorResults = append(errorResults, ErrorResult[T]{
				Index:    i,
				Response: response,
				Err:      response.GetError(),
			})
		}
	}

	return len(errorResults) > 0, errorResults
}

// CheckIfResponseHasErrors checks if a response object that contains Results has any errors
func CheckIfResponseHasErrors[T ResultWithError](response interface{ GetResults() []T }) (bool, ErrorResultList[T]) {
	return CheckIfAnyErrors(response.GetResults())
}
