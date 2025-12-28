package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"shoes-project/models/shoemodel"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	success := ""

	switch r.URL.Query().Get("success") {
	case "created":
		success = "Brand created successfully"
	case "updated":
		success = "Brand updated successfully"
	case "deleted":
		success = "Brand deleted successfully"

	}

	shoes := shoemodel.GetAll()
	data := map[string]any{
		"shoes":   shoes,
		"success": success,
	}

	t := template.New("index.html").Funcs(funcMap)
	t = template.Must(t.ParseFiles("views/shoes/index.html"))

	t.Execute(w, data)
}

// Struct to pass data to template
type FormData struct {
	Brands  []entities.Brand
	Errors  []shoemodel.ValidationError
	OldData map[string]string
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/shoes/create.html")
		if err != nil {
			panic(err)
		}

		brands := brandmodel.GetAll()
		data := FormData{
			Brands:  brands,
			Errors:  nil,
			OldData: make(map[string]string),
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var shoe entities.Shoe

		// Get form values and store for repopulation
		oldData := make(map[string]string)
		oldData["name"] = r.FormValue("name")
		oldData["type"] = r.FormValue("type")
		oldData["description"] = r.FormValue("description")
		oldData["sku"] = r.FormValue("sku")
		oldData["price"] = r.FormValue("price")
		oldData["stock"] = r.FormValue("stock")
		oldData["id_brand"] = r.FormValue("id_brand")

		// Parse brand ID
		brandIdStr := r.FormValue("id_brand")
		if brandIdStr != "" {
			brandId, err := strconv.Atoi(brandIdStr)
			if err == nil {
				shoe.Stock = int64(brandId)
			}
		}

		// Parse stock
		stockStr := r.FormValue("stock")
		if stockStr != "" {
			stock, err := strconv.Atoi(stockStr)
			if err == nil {
				shoe.Stock = int64(stock)
			}
		}

		// Parse price
		priceStr := r.FormValue("price")
		if priceStr != "" {
			stock, err := strconv.Atoi(priceStr)
			if err == nil {
				shoe.Stock = int64(stock)
			}
		}

		// Set other fields

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
