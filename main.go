package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := GetConfig()

	InitApis()
	go CollectTempData()

	fmt.Printf("serving on port %d\n", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
