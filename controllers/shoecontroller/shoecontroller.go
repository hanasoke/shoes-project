package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"shoes-project/models/shoemodel"
	"strconv"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t := template.New("index.html").Funcs(funcMap)
	t = template.Must(t.ParseFiles("views/shoes/index.html"))

	t.Execute(w, nil)
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

		if formData["shoe_price"].(string) == "" {
			validationErrors["shoe_price"] = "Price is required"
		} else {
			price, err := strconv.ParseInt(formData["shoe_price"].(string), 10, 64)
			if err != nil || price <= 0 {
				validationErrors["shoe_price"] = "Price must be a positive number"
			}
		}

		if formData["shoe_stock"].(string) == "" {
			validationErrors["shoe_stock"] = "Stock is required"
		} else {
			stock, err := strconv.ParseInt(formData["shoe_stock"].(string), 10, 64)
			if err != nil || stock < 0 {
				validationErrors["shoe_stock"] = "Stock must be a non-negative number"
			}
		}

		// If there are validation errors, show the form again with errors
		if len(validationErrors) > 0 {
			temp, err := template.ParseFiles("views/shoes/create.html")

			if err != nil {
				panic(err)
			}

			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands":    brands,
				"errors":    validationErrors,
				"form":      formData,
				"duplicate": "please fix the validation errors below",
			}

			temp.Execute(w, data)
			return
		}

		// Convert form values after validation
		brandId, _ := strconv.Atoi(formData["brand_id"].(string))
		stock, _ := strconv.ParseInt(formData["shoe_stock"].(string), 10, 64)
		price, _ := strconv.ParseInt(formData["shoe_price"].(string), 10, 64)

		// Create shoe entity
		shoe.Shoe_Name = formData["shoe_name"].(string)
		shoe.Shoe_Brand.Brand_Id = uint(brandId)
		shoe.Shoe_Type = formData["shoe_type"].(string)
		shoe.Shoe_Description = formData["shoe_description"].(string)
		shoe.Shoe_Sku = formData["shoe_sku"].(string)
		shoe.Shoe_Price = price
		shoe.Shoe_Stock = stock
		shoe.CreatedAt = time.Now()

		// Create shoe in database
		if ok := shoemodel.Create(shoe); !ok {
			// If creation fails, show error
			temp, err := template.ParseFiles("views/shoes/create.html")

			if err != nil {
				panic(err)
			}

			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"error":  "Failed to create shoe. Please try again.",
				"form":   formData,
			}

			temp.Execute(w, data)
			return
		}

		// Success - redirect with success message
		temp, err := template.ParseFiles("views/shoes/create.html")
		if err != nil {
			panic(err)
		}

		brands := brandmodel.GetAll()
		data := map[string]any{
			"brands":  brands,
			"success": "Shoe created successfully!",
			"form":    make(map[string]interface{}), // Clear form
		}

		temp.Execute(w, data)
	}
}
