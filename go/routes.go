package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		HandlerView,
	},
	Route{
		"health",
		"GET",
		"/health",
		Health,
	},
	Route{
		"validator",
		"GET",
		"/validator",
		Validator,
	},
}