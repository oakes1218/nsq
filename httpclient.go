package main

import (
	"log"
	"net/http"
)

func ApiPing(url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8989/"+url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
