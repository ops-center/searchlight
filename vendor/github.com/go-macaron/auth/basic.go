package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"gopkg.in/macaron.v1"
)

// User is the authenticated username that was extracted from the request.
type User string

// BasicRealm is used when setting the WWW-Authenticate response header.
var BasicRealm = "Authorization Required"
var basicPrefix = "Basic "

// Basic returns a Handler that authenticates via Basic Auth. Writes a http.StatusUnauthorized
// if authentication fails.
func Basic(username string, password string) macaron.Handler {
	var siteAuth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		auth := req.Header.Get("Authorization")
		if !SecureCompare(auth, basicPrefix+siteAuth) {
			basicUnauthorized(res)
			return
		}
		c.Map(User(username))
	}
}

// BasicFunc returns a Handler that authenticates via Basic Auth using the provided function.
// The function should return true for a valid username/password combination.
func BasicFunc(authfn func(string, string) bool) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		auth := req.Header.Get("Authorization")
		n := len(basicPrefix)
		if len(auth) < n || auth[:n] != basicPrefix {
			basicUnauthorized(res)
			return
		}
		b, err := base64.StdEncoding.DecodeString(auth[n:])
		if err != nil {
			basicUnauthorized(res)
			return
		}
		tokens := strings.SplitN(string(b), ":", 2)
		if len(tokens) != 2 || !authfn(tokens[0], tokens[1]) {
			basicUnauthorized(res)
			return
		}
		c.Map(User(tokens[0]))
	}
}

func basicUnauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\""+BasicRealm+"\"")
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
