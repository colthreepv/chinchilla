package main

import (
	"encoding/json"
)

type ChiError struct {
	Message string
}

func NewChiError(msg string) string {
	jsonMessage, _ := json.Marshal(ChiError{Message: msg})
	return string(jsonMessage)
}
