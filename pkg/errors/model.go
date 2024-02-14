package errors

import "fmt"

type Error struct {
	Code       string
	Err        error
	Who        string
	StatusHTTP int
	Data       interface{}
	APIMessage string
	UserID     string
}

func NewError() Error {
	return Error{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, Err: %v, Who: %s, Status: %d, Data: %v, UserID: %s",
		e.Code,
		e.Err,
		e.Who,
		e.StatusHTTP,
		e.Data,
		e.UserID,
	)
}
