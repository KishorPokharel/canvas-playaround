package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./canvas-playaround.db")
	if err != nil {
		log.Fatal(err)
	}

	// create table
	createTableSQL := `
    create table if not exists data (
        id integer primary key autoincrement,
        dataurl text
    )
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", handleHome(db))
	http.HandleFunc("/save", handleSaveCanvas(db))

	log.Println("server running")
	http.ListenAndServe(":3000", nil)
}

func handleHome(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            select dataurl from data order by id desc limit 1
        `
		data := ""
		if err := db.QueryRow(query).Scan(&data); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	}
}

type canvasData struct {
	Data string `json:"data"`
}

func handleSaveCanvas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cd := &canvasData{}
		if err := json.NewDecoder(r.Body).Decode(cd); err != nil {
			log.Println(err)
			http.Error(w, "could not decode json", http.StatusBadRequest)
			return
		}
		insertSQL := ` insert into data (dataurl) values ($1) `
		_, err := db.Exec(insertSQL, cd.Data)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(`{"message": "OK"}`)
	}
}
