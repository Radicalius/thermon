package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func parseIntQueryParam(w http.ResponseWriter, r *http.Request, name string) (bool, int64) {
	paramStr := r.URL.Query().Get(name)
	ret, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		fmt.Printf("error parsing query param %s: %s\n", name, err.Error())
		w.WriteHeader(400)
		return false, 0
	}

	return true, ret
}

func InitApis() {
	http.HandleFunc("/api/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		s, from := parseIntQueryParam(w, r, "from")
		if !s {
			return
		}

		s, to := parseIntQueryParam(w, r, "to")
		if !s {
			return
		}

		db := GetDatabase()

		sensorReadings, err := db.Get(from, to)
		if err != nil {
			fmt.Printf("error reading from database: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}

		data, err := json.Marshal(sensorReadings)
		if err != nil {
			fmt.Printf("error serializing sensor readings: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}

		w.Write(data)
	})
}
