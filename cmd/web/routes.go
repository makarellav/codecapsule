package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippets", app.snippetForm)
	mux.HandleFunc("POST /snippets", app.snippetCreate)
	mux.HandleFunc("GET /snippets/{id}", app.snippetView)

	return mux
}
