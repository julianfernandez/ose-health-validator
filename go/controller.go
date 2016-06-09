package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"os"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Health-Validator Service!\n")
}

func Health(w http.ResponseWriter, r *http.Request) {

	var healthOutput HealthOutput

	healthOutput.Status = "Up"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(healthOutput); err != nil {
		panic(err)
	}
}

//endpoint (GET) to validate project health endpoints
func Validator(w http.ResponseWriter, r *http.Request) {

	var token, apiServer string
	apiServerHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	apiServerPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	projectName := os.Getenv("PROJECT_NAME")

	if apiServerHost == "" {
		apiServer = "https://api.boae.paas.gsnetcloud.corp:8443"
	}else {
		apiServer = "https://" + apiServerHost + ":" + apiServerPort
	}
	
	if projectName == "" {
		projectName = "devstack-dev"
	}
	

	// read the service account secret token file at once
	tokenBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Println("Not Service Account Token available")
		token = "uXQAIbJhQESrSwRzajEbUAorau6rPyyM7GC5o86Y7NE"
	} else {
		token = string(tokenBytes[:])
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}

	// Set up the HTTP request to get Services
	urlGetServices := apiServer + "/api/v1/namespaces/" + projectName + "/services"
	req, err := http.NewRequest("GET", urlGetServices, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		panic(err)
	}
	
	resp, err := cli.Do(req)
	
	if err != nil {
		log.Println("Url Get Services=" + urlGetServices)
		log.Fatal("Error getting Services")
	}
	
	defer resp.Body.Close()

	services, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	serviceArray := make([]ServiceObject, 0)
	servicesCustom := map[string]interface{}{}
	json.Unmarshal(services, &servicesCustom)

	if servicesCustom != nil && len(servicesCustom)> 0 {
		items := servicesCustom["items"].([]interface{})

		for _, item := range items {
			itemObject := item.(map[string]interface{})
			metadataMap := itemObject["metadata"].(map[string]interface{})
			specMap := itemObject["spec"].(map[string]interface{})
			serviceArray = append(serviceArray, ServiceObject{metadataMap["name"].(string), specMap["clusterIP"].(string), ValidateService("Hola")})
		}
	}
	
	var validatorOutput ValidatorOutput
	validatorOutput.Services = serviceArray
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(validatorOutput); err != nil {
		panic(err)
	}

}
