package main

import (
	"errors"
	"fmt"
	"github.com/makarellav/codecapsule/internal/models"
	"github.com/makarellav/codecapsule/internal/validator"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.gohtml", data)
}

func (app *application) snippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 1,
	}

	app.render(w, r, http.StatusOK, "create.gohtml", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)

		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)

		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be empty")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be empty")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "create.gohtml", data)

		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.gohtml", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)

		return
	}

	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)

		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		default:
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your sign up was successful. Please log in.")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, r, http.StatusOK, "signup.gohtml", data)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)

		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)

		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			data := app.newTemplateData(r)
			form.AddNonFieldError("Invalid Credentials")
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		default:
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, r, http.StatusOK, "login.gohtml", data)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
