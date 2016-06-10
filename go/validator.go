package main

import (
	"net/http"
	"log"
)

func ValidateService(cli *http.Client, url string, token string) string {
	
	// Set up the HTTP request to validate service/routes
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := cli.Do(req)

	if err != nil {
		log.Println("Health URL = " + url)
		log.Fatal("URL error")
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return "ok"
	} else {
		return "ko"
	}
	
}