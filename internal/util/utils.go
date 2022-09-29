package util

import (
	"github.com/kafka-reader/internal/domain"
)

func IsFatal(err error) bool {
	if e, ok := err.(*domain.AppError); ok && e.Fatal {
		return true
	}
	return false
}
