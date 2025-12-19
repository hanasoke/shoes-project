package main

import (
	"log"
	"net/http"
	"shoes-project/config"
	"shoes-project/controllers/homepage"
	"shoes-project/controllers/shoes"
)

func main() {
	config.ConnectDB()

	// 1. Homepage
	http.HandleFunc("/", homepage.Index)

	// 2. Brand

	// 3. Shoes
	http.HandleFunc("/shoes", shoes.Index)
	http.HandleFunc("/shoe/add", shoes.Add)

	log.Println("Server running on port 8082")
	http.ListenAndServe(":8082", nil)
}
