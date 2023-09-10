package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.tmgasek.net/internal/models"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

// only parse files once when app starts, then store the parsed templates in
// an in memory cache

func newTemplateCache() (map[string]*template.Template, error) {
	// init new map to act as the cache
	cache := map[string]*template.Template{}

	// get slice of all filepaths that match "./ui/html/pages/*.html"
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract file name ("home.html") from full filepath
		name := filepath.Base(page)

		// parse base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// progressively add more templates to the same template set

		// call Parse glob on this ts to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// call parse files on this ts to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template to map using name of the page as key ("home.html")
		cache[name] = ts
	}

	return cache, nil
}
