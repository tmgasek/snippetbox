package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// write error msg and stack trace to errorLog, send generic 500 res to user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// send specific status code and corresponding description to user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// convenient wrapper around clientError, sends 404 res to user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Easily render templates from the cache
func (app *application) render(
	w http.ResponseWriter,
	status int,
	page string,
	data *templateData,
) {
	// retrieve the right template set from cache based on page name
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// init new buffer
	buf := new(bytes.Buffer)

	// write template to buffer instead of to the http.ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if template written to buffer w/o errors, we are good to go.
	w.WriteHeader(status)

	buf.WriteTo(w)
}
