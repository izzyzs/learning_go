package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"
)

// making this function a method against the application struct allows for the use of the loggers I made
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	// page tmpl has to come first,
	// i.e., the page info before the layout
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
			app.serverError(w, err)
		return
	}

	// execute them. Notice how we are passing in the snippet data
	// s: (a models.Snippet struct) as the final parameter
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))	/* gets the id from url query and converts 
	it to int */
	if err != nil || id < 1 { // checks for error or if id is less than 1
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippet: s}

		// page tmpl has to come first,
	// i.e., the page info before the layout
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	// execute them. Notice how we are passing in the snippet data
	// s: (a models.Snippet struct) as the final parameter
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Whispers through the pines"
	content := "Whispers through the pines,\nautumn leaves in twilight's grace—\ntime's quiet embrace."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 	path := r.URL.Path
// 	cleaned := filepath.Clean(path)
// 	http.ServeFile(w, r, cleaned)
// }