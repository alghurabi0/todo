package main

import (
	"context"
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(("Referrer-Policy"), "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the cookie "auth" is set and verify it
		ctx := context.Background()
		cookie, err := r.Cookie("auth")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		client, err := app.users.Auth.Auth(ctx)
		if err != nil {
			app.serverError(w, err)
			return
		}
		t, err := client.VerifyIDToken(ctx, cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			app.errorLog.Println(err)
			return
		}
		ctx = context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		ctx = context.WithValue(ctx, userIdContextKey, t.UID)
		r = r.WithContext(ctx)
		// set the "Cache-Control" header to "no-store" to prevent caching of pages that require authentication
		// in the user's browser
		// w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) protectRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
		next.ServeHTTP(w, r)
	})
}
