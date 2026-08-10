package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const serviceName = "log-service"

func createLogsDirIfNotExists() error {
	if err := os.Mkdir("logs", 0744); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}

func getOrCreateLogFile() (*os.File, error) {
	t := time.Now()

	f, err := os.OpenFile(
		fmt.Sprint("logs/", t.Format(time.DateOnly), ".txt"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func main() {
	if err := createLogsDirIfNotExists(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/log", logHandler)

	srv := http.Server{
		Addr: ":7070",
	}

	log.Fatal(srv.ListenAndServe())
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	f, err := getOrCreateLogFile()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(f, nil)).With(
		slog.String("Req host", r.Host),
	)

	logger.Info(r.RequestURI)

	switch r.Method {
	case "POST":
		handleLog(logger, w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

func handleLog(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	body := r.Body

	/* From documentation
	The Server will close the request body. The ServeHTTP Handler does not need to.
	TODO: to check if body close defer should be removed
	*/
	defer func() {
		if err := body.Close(); err != nil {
			logger.Error(
				fmt.Sprintf("Error closing request body: %v", err),
				slog.String("service-name", serviceName),
			)
		}
	}()

	b, err := io.ReadAll(body)
	if err != nil {
		logger.Error(
			fmt.Sprintf("Error reading request body: %v:", err),
			slog.String("service-name", serviceName),
		)
	}
	// TODO: check if body is not empty

	logger.Info(string(b))

	w.WriteHeader(http.StatusNoContent)
}
