package auth

import (
	"context"
	"database/sql"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "не авторизирован", http.StatusUnauthorized)
			return
		}

		user, err := GetUserByToken(db, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "не авторизирован", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user.ID)

		next(w, r.WithContext(ctx))
	}
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
