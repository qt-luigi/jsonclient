package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const usage = `JSON send/recieve tester.
[args] url jsonfile
    url         sending JSON URL.  (ex. http://localhost:8080/context/sub/action)
    jsonfile    sending JSON data. (ex. {"userid":"hoge", "passwd":"fuga"})`

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	if _, err := url.Parse(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fi, err := os.Stat(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data, err := os.ReadFile(fi.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !json.Valid(data) {
		fmt.Fprintln(os.Stderr, "invalid JSON")
		os.Exit(1)
	}

	status, body, err := sendJSON(os.Args[1], data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Status %d: %s\n", status, http.StatusText(status))
	fmt.Println()
	fmt.Println(body)
}

func sendJSON(url string, data []byte) (int, string, error) {
	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, body, "", "  "); err != nil {
		return 0, "", err
	}

	return res.StatusCode, buf.String(), nil
}
