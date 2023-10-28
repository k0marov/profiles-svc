package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClientError struct {
	DisplayMessage string `json:"message,omitempty"`
	HTTPCode       int    `json:"-"`
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("an error which will be displayed to the client: %d %v", ce.HTTPCode, ce.DisplayMessage)
}

func WriteErrorResponse(w http.ResponseWriter, e error) {
	if ce, ok := e.(ClientError); ok {
		w.WriteHeader(ce.HTTPCode)
		json.NewEncoder(w).Encode(ce)
	}
	w.WriteHeader(http.StatusInternalServerError)
}
