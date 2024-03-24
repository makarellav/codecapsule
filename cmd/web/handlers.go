package main

import (
	"errors"
	"fmt"
	"github.com/makarellav/codecapsule/internal/models"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		log.Print(err.Error())

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		log.Print(err.Error())

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)

		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w, r)
		default:
			app.serverError(w, r, err)
		}

		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}
