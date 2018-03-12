package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

// authorized middleware function that verifies a jwt token
// to enforce user identification
func authorized(w http.ResponseWriter, r *http.Request, secret []byte) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "This endpoint requires a Bearer token")
		return false
	}
	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		writeError(w, 401, "This endpoint requires a Bearer token")
		return false
	}
	token, err := jwt.Parse(matches[1], func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		writeError(w, 401, err.Error())
		return false
	}
	claims := token.Claims.(jwt.MapClaims)
	_, ok := claims["exp"]
	if !ok {
		writeError(w, 401, "missing expiration (exp) claim")
		return false
	}
	return true
}

//JWTSecure is a high order function that based on a secret returns
//a function wrapper to secure http handler functions
func JWTSecure(secret string) func(func(w http.ResponseWriter, r *http.Request)) http.Handler {
	token := []byte(secret)
	return func(fn func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authorized(w, r, token) {
				fn(w, r)
			}
		})
	}
}
