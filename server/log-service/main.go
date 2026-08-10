package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func getOrCreateLogFile() *os.File {
	logsDirName := "logs"

	if err := os.Mkdir(logsDirName, 0644); err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatal(err)
	}

	if err := os.Chdir(logsDirName); err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	f, err := os.OpenFile(
		// TODO: consider separate log files like for access logs and other logs
		fmt.Sprint(t.Format(time.DateOnly), ".txt"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func writeToFile(f *os.File, data []byte) {
	if _, err := f.Write(data); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
}

func main() {
	f := getOrCreateLogFile()

	// TODO: configure server and add handlers
	writeToFile(f, []byte("hello there\n"))
	writeToFile(f, []byte("how are you?\n"))

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
