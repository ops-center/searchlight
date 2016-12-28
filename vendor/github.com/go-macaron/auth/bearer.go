package auth

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

var bearerPrefix = "Bearer "

// Bearer returns a Handler that authenticates via Bearer Auth. Writes a http.StatusUnauthorized
// if authentication fails.
func Bearer(token string) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		auth := req.Header.Get("Authorization")
		if !SecureCompare(auth, bearerPrefix+token) {
			bearerUnauthorized(res)
			return
		}
		c.Map(User(""))
	}
}

// BearerFunc returns a Handler that authenticates via Bearer Auth using the provided function.
// The function should return true for a valid bearer token.
func BearerFunc(authfn func(string) bool) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		auth := req.Header.Get("Authorization")
		n := len(bearerPrefix)
		if len(auth) < n || auth[:n] != bearerPrefix {
			bearerUnauthorized(res)
			return
		}
		if !authfn(auth[n:]) {
			bearerUnauthorized(res)
			return
		}
		c.Map(User(""))
	}
}

func bearerUnauthorized(res http.ResponseWriter) {
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
