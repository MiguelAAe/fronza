package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"portal-backend/pkg/auth"
	"portal-backend/pkg/usersvc"

	"golang.org/x/crypto/bcrypt"
)

// Auth middleware authenticates all type of users
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, ok := auth.ParseToken(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthCourier middleware authenticates only couriers users
func AuthCourier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claim, ok := auth.ParseToken(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claim.Role != usersvc.Courier {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthCustomer middleware authenticates only customers users
func AuthCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claim, ok := auth.ParseToken(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claim.Role != usersvc.Customer {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthAdmin middleware authenticates only admin users
func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claim, ok := auth.ParseToken(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claim.Role != usersvc.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetClaim(r *http.Request) (auth.CustomClaims, bool) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return auth.CustomClaims{}, false
	}

	claims, ok := auth.ParseToken(token)
	if !ok {
		return auth.CustomClaims{}, false
	}

	return claims, true
}

// ContentType middleware
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := r.Header.Get("Content-Type")
		if content != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// APIKey middleware
func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("API-Key")

		claim, ok := GetClaim(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hashedAPIKey, err := usersvc.GetAPIKey(claim.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(hashedAPIKey, []byte(apiKey))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteErr(w http.ResponseWriter, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)
	err2 := json.NewEncoder(w).Encode(&ErrorResponse{
		Error: err.Error(),
	})
	if err2 != nil {
		log.Printf("err sending response: %v\n", err)
	}
	if message != "" {
		log.Printf("%s with err %v", message, err)
		return
	}
	log.Printf("err %v", err)
}
