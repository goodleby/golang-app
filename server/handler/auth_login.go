package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goodleby/golang-app/client"
)

type AuthLoginPayload struct {
	Role string `json:"role"`
	Key  string `json:"key"`
}

type TokenCreator interface {
	NewToken(ctx context.Context, role, key string) (token string, expires time.Time, err error)
}

func AuthLogin(tokenCreator TokenCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload AuthLoginPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding auth payload: %v", err), http.StatusBadRequest, false)
			return
		}

		token, expires, err := tokenCreator.NewToken(ctx, payload.Role, payload.Key)
		if err != nil {
			switch err.(type) {
			case client.ErrUnauthorized:
				HandleError(ctx, w, fmt.Errorf("error creating new auth token: %v", err), http.StatusUnauthorized, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error creating new auth token: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusNoContent)
	}
}
