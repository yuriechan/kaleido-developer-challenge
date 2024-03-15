package middleware

import (
	"backend/internal/utils"
	"net/http"
)

func SetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := r.Header.Get("UserID")
		if len(uid) == 0 {
			// For now, when no user ID header is set, treat as unauthorized request
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := utils.NewContext(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
