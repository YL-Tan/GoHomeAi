package httpresponse

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, isSuccess bool, message string, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: isSuccess,
		Message: message,
		Data:    data,
	}
	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}
