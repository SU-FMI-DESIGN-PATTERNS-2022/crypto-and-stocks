package middleware

import (
	"net/http"
)

func SetContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json")
		next.ServeHTTP(res, req)
	}
}
