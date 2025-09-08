package security

import (
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type BasicAuthChecker struct{}

func NewBasicAuthChecker() *BasicAuthChecker {
	return &BasicAuthChecker{}
}

func (b *BasicAuthChecker) Middleware(decorated http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if check(request) {
			decorated.ServeHTTP(response, request)
			return
		}

		response.WriteHeader(http.StatusUnauthorized)
	}
}

func check(request *http.Request) bool {
	username, password, ok := request.BasicAuth()

	if !ok {
		return false
	}

	if username != os.Getenv("API_USERNAME") {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(os.Getenv("API_PASSWORD")), []byte(password))

	return err == nil
}
