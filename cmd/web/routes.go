package main

import (
	"github.com/makarellav/codecapsule/ui"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.home)))))

	mux.Handle("GET /snippets/{id}", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.snippetView)))))
	mux.Handle("GET /snippets/create", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(app.requireAuth(http.HandlerFunc(app.snippetForm))))))
	mux.Handle("POST /snippets", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(app.requireAuth(http.HandlerFunc(app.snippetCreate))))))

	mux.Handle("GET /users/login", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.userLoginForm)))))
	mux.Handle("POST /users/login", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.userLogin)))))

	mux.Handle("GET /users/signup", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.userSignupForm)))))
	mux.Handle("POST /users/signup", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(http.HandlerFunc(app.userSignup)))))

	mux.Handle("POST /users/logout", app.sessionManager.LoadAndSave(app.csrf(app.authenticate(app.requireAuth(http.HandlerFunc(app.userLogout))))))

	return app.recoverer(app.logRequest(commonHeaders(mux)))
}
