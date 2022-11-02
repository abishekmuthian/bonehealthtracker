// Package secure adds headers to protect against xss and reflection attacks and force use of https
package secure

import (
	"net/http"
)

// These package level variables should be called if required to set policies before the middleware is added

// ContentSecurityPolicy defaults to a strict policy disallowing iframes and scripts from any other origin save self (and Google Analytics for scripts)
var ContentSecurityPolicy = "frame-ancestors 'self'; connect-src 'self'; frame-src 'self' challenges.cloudflare.com github.com/sponsors/abishekmuthian/card; style-src 'self' 'unsafe-inline'; script-src 'self' challenges.cloudflare.com; img-src 'self'"

// Middleware adds some headers suitable for secure sites
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Nonce Implementation Test
		/*		if r.RequestURI=="/projects/110/subscriptions?redirecturl=%2fprojects%2f110-testing-url-validation-2-by-http-status-http-status"{

				var nonce, _ = auth.NonceToken(w,r)

				fmt.Println("Nonce: ",nonce)

				var ContentSecurityPolicy = "frame-ancestors 'self'"+" 'nonce-"+nonce+"'"+"; style-src 'self'"+" 'nonce-"+nonce+"'"+" *.googleapis.com *.paypal.com; script-src 'self'"+" 'nonce-"+nonce+"'"+" www.googletagmanager.com www.google-analytics.com *.paypal.com"
			}*/

		// Add some headers for security

		// Allow no iframing - could also restrict scripts to this domain only (+GA?)
		w.Header().Set("Content-Security-Policy", ContentSecurityPolicy)

		// Allow only https connections for the next 2 years, requesting to be preloaded
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Set ReferrerPolicy explicitly to send only the domain, not the path
		w.Header().Set("Referrer-Policy", "strict-origin")

		// Ask browsers to block xss by default
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Don't allow browser sniffing for content types
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Call the handler
		h(w, r)

	}
}

// HSTSMiddleware adds only the Strict-Transport-Security with a duration of 2 years
func HSTSMiddleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Allow only https connections for the next 2 years, requesting to be preloaded
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Call the handler
		h(w, r)

	}
}
