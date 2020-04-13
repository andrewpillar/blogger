package web

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, msg string, statusCode int) {
	JSON(w, map[string]string{"message": msg}, statusCode)
}
