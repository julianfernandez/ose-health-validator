package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
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

func ValidateReplicas(cli *http.Client, url string, token string) float64  {
	
	replicas :=  0.0
	
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := cli.Do(req)

	if err != nil {
		log.Println("error getDC = " + url)
		return replicas
	}

	defer resp.Body.Close()
	
	dcs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	dcsCustom := map[string]interface{}{}
	json.Unmarshal(dcs, &dcsCustom)

	if dcsCustom != nil && len(dcsCustom)> 0 {
		specMap := dcsCustom["spec"].(map[string]interface{})
		replicas =  specMap["replicas"].(float64)
	}
	
	return replicas
}
