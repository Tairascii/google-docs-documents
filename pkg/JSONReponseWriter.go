package pkg

import (
	"encoding/json"
	"net/http"
)

func JSONErrorResponseWriter(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	data := map[string]string{
		"message": err.Error(),
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func JSONResponseWriter[T any](w http.ResponseWriter, data T, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	dataWrapper := map[string]T{
		"data": data,
	}
	if err := json.NewEncoder(w).Encode(dataWrapper); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
