package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	internalServerErrorMessage = "Internal Server Error"
)

type errorMessage struct {
	Message string `json:"message"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func writeJSON(rw http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func serverError(rw http.ResponseWriter, err error) {
	errorMessage := errorMessage{Message: internalServerErrorMessage}
	werr := writeJSON(rw, http.StatusInternalServerError, errorMessage)
	if werr != nil {
		slog.Error("Error writing error message: " + werr.Error())
	}
	slog.Error("Internal error server: " + err.Error())
}
