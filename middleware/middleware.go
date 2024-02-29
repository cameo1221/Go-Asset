package middleware 

import (
	"net/http"
	"log"
	"io"
	"bytes"

)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)

	})
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("REQUEST:", r)
		// Read request body into bytes
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		// Restore the io.ReadCloser to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Println("REQUEST BODY:", string(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
