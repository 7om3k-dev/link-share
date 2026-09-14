package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "user-links"

var (
	db *pgxpool.Pool
)

type UserLink struct {
	Id          int64
	Title       string
	Description *string
	Url         string
}

func main() {
	dbUserPassword := os.Getenv("DB_PASSWORD")
	if dbUserPassword == "" {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: password not set or empty")
	}

	databaseUrl := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("postgres", dbUserPassword),
		Host:   "user-links-db:5432",
		Path:   "postgres",
	}

	poolConfig, err := pgxpool.ParseConfig(databaseUrl.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse database url: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
	}

	http.HandleFunc("/api/user-links" /* "localhost/api/user-link/{id}" */, func(w http.ResponseWriter, r *http.Request) {
		userLinkHandler(ctx, w, r)
	})

	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":5001", nil))
}

func userLinkHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUserLinks(ctx, w, r)

	//case "PUT":
	//	putUrlHandler(w, r)
	//
	//case "DELETE":
	//	deleteUrlHandler(w, r)

	default:
		w.Header().Add("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getUserLinks(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query(ctx, "SELECT id, title, description, url FROM link")
	if err != nil {
		http.Error(w, "Cannot get user links", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Query error: %v", err)
		return
	}
	defer rows.Close()

	var links []UserLink
	for rows.Next() {
		var link UserLink

		if err := rows.Scan(&link.Id, &link.Title, &link.Description, &link.Url); err != nil {
			http.Error(w, "Cannot serialize user links", http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Row scan error: %v", err)
			return
		}

		links = append(links, link)
	}

	if rows.Err() != nil {
		http.Error(w, "Something went wrong while getting user links", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Rows error: %v", rows.Err())
		return
	}

	if links == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(links); err != nil {
		http.Error(w, "Cannot send json with user links", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "JSON encoding error: %v", err)
		return
	}
}
