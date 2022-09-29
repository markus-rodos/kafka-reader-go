package domain

import (
	"fmt"
)

type AppError struct {
	Fatal   bool
	Message string
	Payload string
	Cause   error
}

func (r *AppError) Error() string {
	return fmt.Sprintf("fatal: %v, message: %v, payload: %v, cause: %v", r.Fatal, r.Message, r.Payload, r.Cause)
}
