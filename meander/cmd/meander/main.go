package main

import (
	"encoding/json"
	"github.com/avalchev94/go_blueprints/meander"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	meander.APIKey = "AIzaSyC6zE-AQjJBAhNXysAcT1wpPMgiP6MTFBA"
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))
	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		urlValues := r.URL.Query()

		q := &meander.Query{
			Journey: strings.Split(urlValues.Get("journey"), "|"),
		}

		var err error
		q.Lat, err = strconv.ParseFloat(urlValues.Get("lat"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		q.Lng, err = strconv.ParseFloat(urlValues.Get("lng"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		q.Radius, err = strconv.Atoi(urlValues.Get("radius"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		q.CostRangeStr = urlValues.Get("cost")

		places := q.Run()
		respond(w, r, places)
	}))

	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}
