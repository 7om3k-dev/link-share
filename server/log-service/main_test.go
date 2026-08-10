package main_test

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

/*
Note:

	By documentation, The Client.Transport typically has internal state (cached TCP connections),
	so Clients should be reused instead of created as needed. Clients are safe for concurrent use by multiple goroutines.
*/
var client = &http.Client{}

func createHttpPostRequest(t *testing.T, body *bytes.Buffer) *http.Request {
	// TODO: candidate to be more reusable function
	r, err := http.NewRequest("POST", "http://localhost:7070/log", body)

	if err != nil {
		t.Errorf("TestLogEntry: http request creation error: %v", err)
	}

	r.Host = "example.com:7070"

	return r
}

func readLastLine(scanner *bufio.Scanner) (string, error) {
	var lastLine string

	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

func TestLogEntrySucceeded(t *testing.T) {
	buf := bytes.Buffer{}
	buf.WriteString("Authentication failed")

	r := createHttpPostRequest(t, &buf)

	// By documentation, on error, any Response can be ignored. No need to close it's body.
	resp, err := client.Do(r)
	/*
	  By documentation:
	    - Any returned error will be of type *url.Error.
	    - The url.Error value's Timeout method will report true if the request timed out.
	*/
	if err != nil {
		if uErr, ok := errors.AsType[*url.Error](err); ok && uErr.Timeout() {
			t.Errorf("TestLogEntry: http request timeout")
		}
		t.Errorf("TestLogEntry: http response error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("TestLogEntry: bad response status code: %v", resp.StatusCode)
	}

	tn := time.Now()
	f, err := os.Open("logs/" + tn.Format(time.DateOnly) + ".txt")
	if err != nil {
		t.Errorf("TestLogEntry: error opening file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lastLine, readErr := readLastLine(scanner)

	if readErr != nil {
		t.Errorf("TestLogEntry: error reading log message sent in logs file: %v", readErr)
	}

	hasExpectedEntries := strings.Contains(lastLine, "msg=\"Authentication failed\"") &&
		strings.Contains(lastLine, "\"Req host\"=example.com:7070")

	if !hasExpectedEntries {
		t.Errorf("TestLogEntry: could not find log message in logs: %v", lastLine)
	}
}
