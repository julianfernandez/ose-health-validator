package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	var serviceArray []ServiceObject
	isOSE := false
	apiServerHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	apiServerPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	projectName := os.Getenv("PROJECT_NAME")

	if apiServerHost == "" {
		apiServer = "https://api.boae.paas.gsnetcloud.corp:8443"
	} else {
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
		isOSE = true
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}
	
	// Set up the HTTP request to get Routes
	urlGetRoutes := apiServer + "/oapi/v1/namespaces/" + projectName + "/routes"
	req, err := http.NewRequest("GET", urlGetRoutes, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		panic(err)
	}

	resp, err := cli.Do(req)

	if err != nil {
		log.Println("Url Get Routes=" + urlGetRoutes)
		log.Fatal("Error getting Routes")
	}

	defer resp.Body.Close()

	routes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	routeArray := make([]ServiceObject, 0)
	routesCustom := map[string]interface{}{}
	json.Unmarshal(routes, &routesCustom)

	if routesCustom != nil && len(routesCustom) > 0 {
		items := routesCustom["items"].([]interface{})

		for _, item := range items {
			itemObject := item.(map[string]interface{})
			metadataMap := itemObject["metadata"].(map[string]interface{})
			specMap := itemObject["spec"].(map[string]interface{})
			tls := specMap["tls"]
			var protocol string
			if tls != nil {
				protocol = "https"
			} else {
				protocol = "http"
			}
			urlRoute := protocol + "://" + specMap["host"].(string) + "/health"
			routeArray = append(routeArray, ServiceObject{metadataMap["name"].(string), urlRoute, ValidateService(cli, urlRoute, token)})
		}
	}	

	//If Docker image in running into Openshift
	if isOSE {
		// Set up the HTTP request to get Services
		urlGetServices := apiServer + "/api/v1/namespaces/" + projectName + "/services"
		req, err = http.NewRequest("GET", urlGetServices, nil)
		req.Header.Add("Authorization", "Bearer "+token)

		if err != nil {
			panic(err)
		}

		resp, err = cli.Do(req)

		if err != nil {
			log.Println("Url Get Services=" + urlGetServices)
			log.Fatal("Error getting Services")
		}

		defer resp.Body.Close()

		services, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		serviceArray = make([]ServiceObject, 0)
		servicesCustom := map[string]interface{}{}
		json.Unmarshal(services, &servicesCustom)

		if servicesCustom != nil && len(servicesCustom) > 0 {
			items := servicesCustom["items"].([]interface{})

			for _, item := range items {
				itemObject := item.(map[string]interface{})
				metadataMap := itemObject["metadata"].(map[string]interface{})
				specMap := itemObject["spec"].(map[string]interface{})
				portMap := specMap["ports"].([]interface{})
				portObject := portMap[0].(map[string]interface{})
				urlService := "http://" + specMap["clusterIP"].(string) + ":" + strconv.FormatFloat(portObject["port"].(float64), 'f', 0, 64) + "/health"
				serviceArray = append(serviceArray, ServiceObject{metadataMap["name"].(string), urlService, ValidateService(cli, urlService, token)})
			}
		}
	}

	var validatorOutput ValidatorOutput
	if isOSE{
		validatorOutput.Services = serviceArray
	}
	validatorOutput.Routes = routeArray

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(validatorOutput); err != nil {
		panic(err)
	}

}
