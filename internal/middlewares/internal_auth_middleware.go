/**

 filename  : internal_auth_middleware.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-28

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
)

func InternalAuthMiddleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])
		expectedHash := os.Getenv("INTERNAL_API_TOKEN_HASH")

		if expectedHash == "" {
			http.Error(w, "Server misconfigured", http.StatusInternalServerError)
			return
		}

		if subtle.ConstantTimeCompare([]byte(tokenHash), []byte(expectedHash)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}