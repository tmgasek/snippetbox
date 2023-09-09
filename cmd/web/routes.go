package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// serves files out of .ui/static dir
	// path is relative to projuct root dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register file server as handler for all url paths that start with /static/
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
