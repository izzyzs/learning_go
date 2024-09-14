package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux() // ServeMux == Go's built-in version of a Router
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	// routes different paths to different server handlers defined in handlers.go
	// they are passed into the function bc they are in the same package
	
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// creates an http file server which serves files from a certain directory
	
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// serves said files to requester by first stripping the leading "/static" from the url path  
	// and then passing it into the fileServer

	// mux.HandleFunc("/static/download/", downloadHandler)
	return mux
}