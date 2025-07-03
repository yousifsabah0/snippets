package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCace[page]
	if !ok {
		err := fmt.Errorf("template %s does not exixts!", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	if err := ts.ExecuteTemplate(w, "index", data); err != nil {
		app.serverError(w, r, err)
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, v any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := app.formDecoder.Decode(v, r.PostForm); err != nil {
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.session.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
