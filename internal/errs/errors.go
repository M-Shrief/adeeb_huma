package errs // named like that so it doesn't conflict with the original errors pacakge in golang

import "fmt"

type APIError struct {
	Status   int
	Message  string
	CausedIn string
	Err      error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error, Status: %s", e.Status)
}

func (e *APIError) Unwrap() error {
	return e.Err
}
