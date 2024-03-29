package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		writer.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "deny")
		writer.Header().Set("X-XSS-Protection", "0")

		// Move to next handler if successful
		next.ServeHTTP(writer, request)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())

		next.ServeHTTP(writer, request)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			// Use recover() to check if a panic occured
			if err := recover(); err != nil {
				writer.Header().Set("Connection", "close")

				// send the error as Internal Server Error to the user
				app.serverError(writer, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !app.Authenticate(request) {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}

		writer.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(writer, request)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := app.sessionManager.GetInt(request.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(writer, request)
			return
		}

		// Check to see if a user with that ID exists in our database
		exists, err := app.users.Exists(id)
		if err != nil {
			next.ServeHTTP(writer, request)
			// app.serverError(writer, err)
			return
		}

		if exists {
			ctx := context.WithValue(request.Context(), isAuthenticatedCtxKey, true)
			request = request.WithContext(ctx)
		}

		next.ServeHTTP(writer, request)
	})
}
