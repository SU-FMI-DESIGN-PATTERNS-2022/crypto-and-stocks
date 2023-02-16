package middleware

import (
	"context"
	"fmt"
	"net/http"
)

type UserIdKey string

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		//TODO: add logic for authentication
		var userId int64
		key := UserIdKey("userId")

		ctxWithUserId := context.WithValue(req.Context(), key, userId)
		reqWithUserId := req.WithContext(ctxWithUserId)

		fmt.Println("Authenticating...")

		next.ServeHTTP(res, reqWithUserId)
	}
}
