package util

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string, err error) {
	resp := APIResponse{
		Success: false,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	WriteJSON(w, status, resp)
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, status, resp)
}
