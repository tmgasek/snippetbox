package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// serves files out of .ui/static dir
	// path is relative to projuct root dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register file server as handler for all url paths that start with /static/
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// create a new middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
