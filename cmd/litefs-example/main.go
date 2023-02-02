package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/mattn/go-sqlite3"
)

// Command line flags.
var (
	dsn  = flag.String("dsn", "", "datasource name")
	addr = flag.String("addr", ":8080", "bind address")
)

var db *sql.DB

//go:embed schema.sql
var schemaSQL string

func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) (err error) {
	flag.Parse()

	if *dsn == "" {
		return fmt.Errorf("dsn required")
	} else if *addr == "" {
		return fmt.Errorf("bind address required")
	}

	// Connect to SQLite database.
	db, err = sql.Open("sqlite3", *dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	log.Printf("database opened at %s", *dsn)

	// Run migration.
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("cannot migrate schema: %w", err)
	}

	// Start HTTP server.
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/generate", handleGenerate)

	log.Printf("http server listening on %s", *addr)
	return http.ListenAndServe(*addr, nil)
}

//go:embed index.tmpl
var indexTmplContent string
var indexTmpl = template.Must(template.New("index").Parse(indexTmplContent))

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// If a different region is specified, redirect to that region.
	if region := r.URL.Query().Get("region"); region != "" && region != os.Getenv("FLY_REGION") {
		log.Printf("redirecting from %q to %q", os.Getenv("FLY_REGION"), region)
		w.Header().Set("fly-replay", "region="+region)
		return
	}

	// Query for the most recently added people.
	rows, err := db.Query(`
		SELECT id, name, phone, company
		FROM persons
		ORDER BY id DESC
		LIMIT 10
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect rows into a slice.
	var persons []*Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.ID, &person.Name, &person.Phone, &person.Company); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		persons = append(persons, &person)
	}
	if err := rows.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the list to either text or HTML.
	tmplData := TemplateData{
		Region:  os.Getenv("FLY_REGION"),
		Persons: persons,
	}

	switch r.Header.Get("accept") {
	case "text/plain":
		fmt.Fprintf(w, "REGION: %s\n\n", tmplData.Region)
		for _, person := range tmplData.Persons {
			fmt.Fprintf(w, "- %s @ %s (%s)\n", person.Name, person.Company, person.Phone)
		}

	default:
		if err := indexTmpl.ExecuteTemplate(w, "index", tmplData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	// Only allow POST methods.
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/*
		// If this node is not primary, look up and redirect to the current primary.
		primaryFilename := filepath.Join(filepath.Dir(*dsn), ".primary")
		primary, err := os.ReadFile(primaryFilename)
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if string(primary) != "" {
			log.Printf("redirecting to primary instance: %q", string(primary))
			w.Header().Set("fly-replay", "instance="+string(primary))
			return
		}
	*/

	// If this is the primary, attempt to write a record to the database.
	person := Person{
		Name:    gofakeit.Name(),
		Phone:   gofakeit.Phone(),
		Company: gofakeit.Company(),
	}
	if _, err := db.ExecContext(r.Context(), `INSERT INTO persons (name, phone, company) VALUES (?, ?, ?)`, person.Name, person.Phone, person.Company); err != nil {
		http.Error(w, "Method not alllowed", http.StatusMethodNotAllowed)
		return
	}

	// Redirect back to the index page to view the new result.
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

type TemplateData struct {
	Region  string
	Persons []*Person
}

type Person struct {
	ID      int
	Name    string
	Phone   string
	Company string
}
