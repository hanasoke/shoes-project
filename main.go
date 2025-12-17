package main

import (
	"log"
	"net/http"
	"shoes-project/config"
	"shoes-project/controllers/brand"
	"shoes-project/controllers/homepage"
)

func main() {
	config.ConnectDB()

	// 1. Homepage
	http.HandleFunc("/", homepage.Index)

	// 2. Brand
	http.HandleFunc("/brand", brand.Index)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
