package main

import (
	"log"
	"net/http"
	"github.com/marove2000/hack-and-pay/router/v1"
)

func main() {

	router := v1.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}