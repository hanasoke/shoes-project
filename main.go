package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
