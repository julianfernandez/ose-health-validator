package main

import (
	"net/http"
	"log"
)

func ValidateService(cli *http.Client, url string, token string) string {
	
	// Set up the HTTP request to validate service/routes
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := cli.Do(req)

	if err != nil {
		log.Println("error Health URL = " + url)
		return "ko"
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 ||  resp.StatusCode == 401 ||resp.StatusCode == 403 {
		return "ok"
	} else {
		return "ko"
	}
	
}