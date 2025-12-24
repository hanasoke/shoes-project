package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"shoes-project/models/shoemodel"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	shoes := shoemodel.GetAll()
	data := map[string]any{
		"shoes": shoes,
	}

	temp, err := template.ParseFiles("views/shoes/index.html")

	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/shoes/create.html")
		if err != nil {
			panic(err)
		}

		brands := brandmodel.GetAll()
		data := map[string]any{
			"brands": brands,
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var shoe entities.Shoe
		var validationErrors = make(map[string]string)
		var formData = make(map[string]interface{})

		// Store form data for re-population
		formData["shoe_name"] = strings.TrimSpace(r.FormValue("shoe_name"))
		formData["brand_id"] = r.FormValue("brand_id")
		formData["shoe_type"] = strings.TrimSpace(r.FormValue("shoe_type"))
		formData["shoe_description"] = strings.TrimSpace(r.FormValue("shoe_description"))
		formData["shoe_sku"] = strings.TrimSpace(r.FormValue("shoe_sku"))
		formData["shoe_price"] = r.FormValue("shoe_price")
		formData["shoe_stock"] = r.FormValue("shoe_stock")

		// Validate required fields
		if formData["shoe_name"].(string) == "" {
			validationErrors["shoe_name"] = "Shoe name is required"
		} else if len(formData["shoe_name"].(string)) > 100 {
			validationErrors["shoe_name"] = "Shoe name cannot exceed 100 characters"
		}

		if formData["brand_id"].(string) == "" {
			validationErrors["brand_id"] = "Brand is required"
		}

		if formData["shoe_type"].(string) == "" {
			validationErrors["shoe_type"] = "Shoe type is required"
		}

		if formData["shoe_sku"].(string) == "" {
			validationErrors["shoe_sku"] = "SKU is required"
		} else {
			// Check for duplicate SKU
			if shoemodel.IsSkuExists(formData["shoe_sku"].(string)) {
				validationErrors["shoe_sku"] = "SKU already exists"
			}
		}

		if ok := shoemodel.Create(shoe); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
