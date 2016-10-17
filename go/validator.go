package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"crypto/tls"
)

func ValidateService(cli *http.Client, url string, token string) string {
	
	// Set up the HTTP request to validate service/routes
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Add("Authorization", "Bearer "+token)
	resp, err := cli.Do(req)

	if err != nil {
		log.Println("error Health URL = " + url)
		return "ko"
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 401 ||resp.StatusCode == 403 {
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
		spec := dcsCustom["spec"]
		if spec != nil {
			specMap := spec.(map[string]interface{})
			replicas =  specMap["replicas"].(float64)
		}
	}
	
	return replicas
}



func ProjectChecker() []ServiceObject  {
	var token, apiServer string
	//	serviceArray := make([]ServiceObject, 0)
	routeArray := make([]ServiceObject, 0)
	//isOSE := false
	apiServerHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	apiServerPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	projectName := os.Getenv("PROJECT_NAME")

	if apiServerHost != "" {
		apiServer = "https://" + apiServerHost + ":" + apiServerPort
	}

	// read the service account secret token file at once
	tokenBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Println("Not Service Account Token available")
	} else {
		token = string(tokenBytes[:])
		//isOSE = true
	}
	apiServer = "https://api.boaw.paas.gsnetcloud.corp:8443"
	projectName = "globalpaas-dev"
	token = "snjiYBqmIo6l2RmkiNTeovFmbwcgxSLu7SP78-oQ2f4"

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
			urlDC :=  apiServer + "/oapi/v1/namespaces/" + projectName + "/deploymentconfigs/" + metadataMap["name"].(string)
			routeArray = append(routeArray, ServiceObject{metadataMap["name"].(string), urlRoute, ValidateService(cli, urlRoute, token), ValidateReplicas(cli, urlDC, token)})
		}
	}

	//	//If Docker image in running into Openshift
	//	if isOSE {
	//		// Set up the HTTP request to get Services
	//		urlGetServices := apiServer + "/api/v1/namespaces/" + projectName + "/services"
	//		req, err = http.NewRequest("GET", urlGetServices, nil)
	//		req.Header.Add("Authorization", "Bearer "+token)
	//
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		resp, err = cli.Do(req)
	//
	//		if err != nil {
	//			log.Println("Url Get Services=" + urlGetServices)
	//			log.Fatal("Error getting Services")
	//		}
	//
	//		defer resp.Body.Close()
	//
	//		services, err := ioutil.ReadAll(resp.Body)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		servicesCustom := map[string]interface{}{}
	//		json.Unmarshal(services, &servicesCustom)
	//
	//		if servicesCustom != nil && len(servicesCustom) > 0 {
	//			items := servicesCustom["items"].([]interface{})
	//
	//			for _, item := range items {
	//				itemObject := item.(map[string]interface{})
	//				metadataMap := itemObject["metadata"].(map[string]interface{})
	//				specMap := itemObject["spec"].(map[string]interface{})
	//				portMap := specMap["ports"].([]interface{})
	//				portObject := portMap[0].(map[string]interface{})
	//				urlService := "http://" + specMap["clusterIP"].(string) + ":" + strconv.FormatFloat(portObject["port"].(float64), 'f', 0, 64) + "/health"
	//				serviceArray = append(serviceArray, ServiceObject{metadataMap["name"].(string), urlService, ValidateService(cli, urlService, token)})
	//			}
	//		}
	//	}

	return routeArray
}
