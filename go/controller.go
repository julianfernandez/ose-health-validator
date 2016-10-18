package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"log"
	//"strconv"
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

	validatorOutput := ProjectChecker()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(validatorOutput); err != nil {
		panic(err)
	}

}

func HandlerView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/views/index.html")
    if err != nil {
        log.Print("template parsing error: ", err)
    }
    var routes []ServiceObject
    routes = ProjectChecker()
    err = t.ExecuteTemplate(w, "index", routes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}    
}
