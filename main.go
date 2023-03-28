package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// create client to send requests
	client := &http.Client{}

	// create request to be signed
	req, err := http.NewRequest("GET", "https://st.letterboxd.com/api/v0/film/lScm", nil)
	if err != nil {
		log.Fatal(err)
	}

	// sign it!
	signRequest(req)

	// send it!
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// read it!
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
