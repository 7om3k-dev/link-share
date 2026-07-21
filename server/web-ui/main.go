package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/7om3k/link-share/applib/logger"
)

const serviceName = "web-ui"

var (
	serviceLogger *logger.Logger
)

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	t, _ := template.ParseFiles("pages/home.html")
	t.Execute(w, nil)
}

func signInPageHandler(w http.ResponseWriter, _ *http.Request) {
	t, _ := template.ParseFiles("pages/sign-in.html")
	t.Execute(w, nil)
}

func signUpPageHandler(w http.ResponseWriter, _ *http.Request) {
	t, _ := template.ParseFiles("pages/sign-up.html")
	t.Execute(w, nil)
}

func userLinksPageHandler(w http.ResponseWriter, _ *http.Request) {
	resp, err := http.Get("http://user-links:5001/api/user-links")

	if err != nil {
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		serviceLogger.LogError(logger.MessageKey, "User links fetch error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		temp, _ := template.ParseFiles("pages/user-links.html")
		temp.Execute(w, nil)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		serviceLogger.LogError(logger.MessageKey, "User links fetch error, bad status: %v\n", resp.StatusCode, err)
		return
	}

	type userLink struct {
		Id          int64
		Title       string
		Description *string
		Url         string
	}
	userLinkList := make([]userLink, 0)

	// TODO: switch to decoder v2, when no longer experimental
	dec := json.NewDecoder(resp.Body)

	// Endpoint returns array of user links, token bellow read open bracket "[" from json.
	_, err = dec.Token()
	if err != nil {
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		serviceLogger.LogError(logger.MessageKey, "User links decode error: "+err.Error())
		return
	}

	for dec.More() {
		var l userLink
		if err := dec.Decode(&l); err != nil {
			http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
			serviceLogger.LogError(logger.MessageKey, "User links decode error: "+err.Error())
			return
		}
		userLinkList = append(userLinkList, l)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		serviceLogger.LogError(logger.MessageKey, "User links decode error: "+err.Error())
		return
	}

	temp, _ := template.ParseFiles("pages/user-links.html")
	temp.Execute(w, userLinkList)
}

func main() {
	serviceLogger = logger.NewAppLogger(serviceName)
	mux := http.NewServeMux()

	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/sign-in", signInPageHandler)
	mux.HandleFunc("/sign-up", signUpPageHandler)
	mux.HandleFunc("/user-links", userLinksPageHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: http.NewCrossOriginProtection().Handler(mux),
	}

	serviceLogger.LogInfo(logger.MessageKey, "Starting web ui service")
	log.Fatal(srv.ListenAndServe())
}
