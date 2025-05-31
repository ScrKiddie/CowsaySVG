package main

import (
	"cowsaysvg/api/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Println("server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
