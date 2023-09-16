package main

import (
	"fmt"
	"net/http"
)

// secureHeaders -> servemux -> app handler
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf(
			"%s - %s %s %s",
			r.RemoteAddr,
			r.Proto,
			r.Method,
			r.URL.RequestURI(),
		)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a deferred func (which will always run in the event of a
		// panic as Go unwinds the stack)
		defer func() {
			// check if there has been a panic
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user is not authed, redirect and return from middleware chain so
		// that no subsequent handlers in the chain are executed.

		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set the cache-control header so that pages that need auth are not
		// stored in the users browser cache (or any other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// Call next handler in chain.
		next.ServeHTTP(w, r)
	})
}
