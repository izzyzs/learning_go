package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {			// To allow our handler functions to use our errorLog and infoLog
	errorLog *log.Logger			// we can make them global variables and make our handler functions
	infoLog *log.Logger				// methods AGAINST this new application struct.
	snippets *mysql.SnippetModel	// Makes the SnippetModel a global variable/dependency injection
}

func main() {
	// Command line flags (i.e., cl flag)
	addr := flag.String("addr", ":3000", "HTTP network address")
	// defines a cl flag addr and set value to addr
	// (cl flag, default value, purpose) -> return string *
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	// defines a new cl flag for the MySQL DSN string

	flag.Parse() // parses the flag and assigns it to the variable; IF NOT HERE, PROGRAM WILL USE 
	// DEFAULTED VALUE

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime) // creates a new logger that 
	// defines the Writer, the prefix (being INFO in this case), and a flag (additional 
	// information) to add (date and time)

	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile) // does the 
	// same as above; adds file information regarding name and line number of the error

	db, err := openDB(*dsn) // db connection method "openDB" can be found below, takes dsn from cl and
	// returns db pointer and err

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close() // defer db.Close() call so it closes before end of main()


	
	app := &application{			// this app initializes the application struct and allows for
		errorLog: errorLog,			// the use of the dependencies
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}


	// Route files were originally here, they were removed for decluttering to routes.go
	// the mux are now injected back into main through the app.routes() method used with the
	// Handler field below

	// *******************************************************
	// CONCLUSION OF REFACTORING/Separation of Concern Changes
	// *******************************************************

	/*
		Something that is really important with regards to software design is the principle of
		seperation of concerns. The idea that every action that has to be handle is separated away
		into it's own section. Routing has it's area, functions that handle requests and responses
		have their own area, even middleware/wrapper functions have their own area. Keep in mind 
		this principle even if working on another project or in another language.
	*/
	

	
	// To get Go's HTTP server to use our new errorLog, we need to initialized an http.Server 
	// struct, it will use same address && routes as before but we will give it our errorLog for 
	// logging errors
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	
	infoLog.Printf("Starting server on PORT %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
	// connects to localhost port, logs to console the attempt to start and any err if it occurs
}

func openDB(dsn string) (*sql.DB, error)  {			// sql.Open() wrapper
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}