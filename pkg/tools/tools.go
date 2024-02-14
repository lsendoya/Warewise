package tools

import (
	"encoding/json"
	)

func UnmarshalJSON(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
