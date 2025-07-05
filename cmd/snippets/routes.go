package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/yousifsabah0/snippets/web"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	middleware := alice.New(app.panicRecovery, app.logRequest, headers)
	dynamic := alice.New(app.session.LoadAndSave, noSurf, app.authenticate)

	mux.HandleFunc("GET /ping", app.ping)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.snippetView))

	mux.Handle("GET /users/signup", dynamic.ThenFunc(app.signupForm))
	mux.Handle("POST /users/signup", dynamic.ThenFunc(app.signup))

	mux.Handle("GET /users/login", dynamic.ThenFunc(app.loginForm))
	mux.Handle("POST /users/login", dynamic.ThenFunc(app.login))

	protected := dynamic.Append(app.requiredAuth)

	mux.Handle("GET /snippets/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippets/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /users/logout", protected.ThenFunc(app.logout))

	mux.Handle("GET /static/", http.FileServerFS(web.Files))

	return middleware.Then(mux)
}
