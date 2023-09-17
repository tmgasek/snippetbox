package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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

// NoSurf middleware func which uses a customised CSRF cookie with the Secure,
// Path, and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get value from session. Will return zero value for int (0) if no
		// such value in session. - Call next handler in chain and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Check if user with that ID exists
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// If exists, we know that the req is coming from authed user in DB.
		// Create a new copy of req with extra ctx val, and assign to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
