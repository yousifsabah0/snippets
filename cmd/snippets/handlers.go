package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yousifsabah0/snippets/internal/models"
	"github.com/yousifsabah0/snippets/internal/validators"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusNotFound)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
			return
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

type snippetCreateForm struct {
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validators.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Title), "title", "This field is required")
	form.CheckField(validators.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validators.NotBlank(form.Content), "content", "This field is required")
	form.CheckField(validators.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equals 1, 7, or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.session.Put(r.Context(), "flash", "A new snippet successfully created")
	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)

	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Wassssssssup. creating a snippet"))
}

type signupForm struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}

	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Name is required")
	form.CheckField(validators.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validators.Matches(form.Email, validators.EmailRx), "email", "yooooo! bad email dude")
	form.CheckField(validators.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "Password must be +8")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	if err := app.users.Insert(form.Name, form.Email, form.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddError("email", "email address is already used")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		}
		app.serverError(w, r, err)
		return
	}

	app.session.Put(r.Context(), "flash", "Your account has been seccessfully created... you can login now.")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type loginForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}

	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Email), "email", "required...")
	form.CheckField(validators.Matches(form.Email, validators.EmailRx), "email", "required to be email...")
	form.CheckField(validators.NotBlank(form.Password), "password", "required...")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Incorrect both")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnsupportedMediaType, "login.html", data)
		}

		app.serverError(w, r, err)
		return
	}

	if err := app.session.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.session.Put(r.Context(), "authID", id)
	http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	if err := app.session.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.session.Remove(r.Context(), "authID")
	app.session.Put(r.Context(), "flash", "You've been logged out successfully.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
