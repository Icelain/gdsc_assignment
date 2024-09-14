package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type StoreRequest struct {
	Text string `json:"text"`
}

type StoreResponse struct {
	TextResp []string `json:"response"`
}

func init() {

	db, _ := sql.Open("sqlite3", "localdb.sql")
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS STORE(data text);")

}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func main() {

	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", "localdb.sql")
	defer db.Close()

	if err != nil {

		log.Fatalf("could not start database: %e", err)

	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {

			http.Error(w, "only get methods allowed on this endpoint", http.StatusBadRequest)
			return

		}

		tmpl := template.Must(template.ParseFiles("./index.html"))
		if err := tmpl.Execute(w, nil); err != nil {

			log.Fatal(err)

		}

	})

	mux.HandleFunc("/uwu", CORS(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {

			http.Error(w, "only get methods allowed on this endpoint", http.StatusBadRequest)
			return

		}

		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{

			"message": "you found the base endpoint uwu ;)",
		})

	}))

	mux.HandleFunc("/api/store_text", CORS(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {

			http.Error(w, "only post methods allowed on this endpoint", http.StatusBadRequest)
			return
		}

		if r.Header.Get("content-type") != "application/json" {

			http.Error(w, "only calls with application/json headers are allowed on this endpoint", http.StatusBadRequest)
			return

		}

		var storeRequest StoreRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&storeRequest); err != nil {

			http.Error(w, "invalid json", http.StatusBadRequest)
			return

		}

		_, err := db.Exec("INSERT INTO STORE(data) VALUES(?);", storeRequest.Text)

		if err != nil {

			log.Println(err)
			http.Error(w, "Something went wrong internally", http.StatusInternalServerError)
			return

		}

		json.NewEncoder(w).Encode(map[string]string{

			"message": "text has been successfully inserted",
		})

	}))

	mux.HandleFunc("/api/get_text", CORS(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {

			http.Error(w, "only get methods allowed on this endpoint", http.StatusBadRequest)
			return

		}

		w.Header().Add("content-type", "application/json")
		rows, err := db.Query("SELECT * FROM STORE;")

		if err != nil {

			http.Error(w, "Sorry there was an internal error", http.StatusInternalServerError)
			return
		}

		var sr StoreResponse
		for rows.Next() {

			var text string
			rows.Scan(&text)
			sr.TextResp = append(sr.TextResp, text)

		}

		if err := json.NewEncoder(w).Encode(sr); err != nil {

			http.Error(w, "request failed internally", http.StatusInternalServerError)
			return

		}

	}))

	log.Println("listening on port 5050")
	http.ListenAndServe(":5050", mux)

}
