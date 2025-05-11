package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meaqese/norpn/internal/orch/services"
	"net/http"
	"strings"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func withAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		header := r.Header.Get("Authorization")

		encoder := json.NewEncoder(w)

		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		id, err := authService.Parse(strings.TrimPrefix(header, prefix))
		if err != nil {
			fmt.Println(err)
			encoder.Encode(&ResponseUser{
				Status: false,
				Error:  err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
