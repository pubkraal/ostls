package util

import (
	"github.com/google/uuid"
)

func ConvertTokensToSet(slice []uuid.UUID) map[string]bool {
	s := make(map[string]bool)
	for k := range slice {
		s[slice[k].String()] = true
	}
	return s
}

func ConvertSliceToSet(slice []interface{}) map[string]bool {
	s := make(map[string]bool)
	for k := range slice {
		s[string(k)] = true
	}
	return s
}
