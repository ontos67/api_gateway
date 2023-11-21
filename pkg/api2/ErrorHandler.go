package api2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func errorHandler(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
	}

	data, err := json.Marshal(response)
	if err != nil {
		// Failed to marshal error response, log the error
		log.Println("Failed to marshal error response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}
