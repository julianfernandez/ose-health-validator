package main

import (

)

type HealthOutput struct {
	Status string `json:"status"`
}

type ServiceObject struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Health string `json:"health"`
	Replicas float64  `json:"replicas"`
}

type ValidatorOutput struct {
//	Services []ServiceObject `json:"services"`
	Routes []ServiceObject `json:"routes"`
}
