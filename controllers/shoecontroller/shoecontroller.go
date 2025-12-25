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
		formData["name"] = strings.TrimSpace(r.FormValue("name"))
		formData["brand_id"] = r.FormValue("brand_id")
		formData["type"] = strings.TrimSpace(r.FormValue("type"))
		formData["description"] = strings.TrimSpace(r.FormValue("description"))
		formData["sku"] = strings.TrimSpace(r.FormValue("sku"))
		formData["price"] = r.FormValue("price")
		formData["stock"] = r.FormValue("stock")

		// Validate required fields
		if formData["name"].(string) == "" {
			validationErrors["name"] = "Shoe name is required"
		} else if len(formData["name"].(string)) > 100 {
			validationErrors["name"] = "Shoe name cannot exceed 100 characters"
		}

		if formData["brand_id"].(string) == "" {
			validationErrors["brand_id"] = "Brand is required"
		}

		if formData["type"].(string) == "" {
			validationErrors["type"] = "Shoe type is required"
		}

		if formData["sku"].(string) == "" {
			validationErrors["sku"] = "SKU is required"
		} else {
			// Check for duplicate SKU
			if shoemodel.IsSkuExists(formData["sku"].(string)) {
				validationErrors["sku"] = "SKU already exists"
			}
		}

		if formData["price"].(string) == "" {
			validationErrors["price"] = "Price is required"
		} else {
			price, err := strconv.ParseInt(formData["price"].(string), 10, 64)
			if err != nil || price <= 0 {
				validationErrors["price"] = "Price must be a positive number"
			}
		}

		if formData["stock"].(string) == "" {
			validationErrors["stock"] = "Stock is required"
		} else {
			stock, err := strconv.ParseInt(formData["stock"].(string), 10, 64)
			if err != nil || stock < 0 {
				validationErrors["stock"] = "Stock must be a non-negative number"
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
		stock, _ := strconv.ParseInt(formData["stock"].(string), 10, 64)
		price, _ := strconv.ParseInt(formData["price"].(string), 10, 64)

		// Create shoe entity
		shoe.Name = formData["name"].(string)
		shoe.Brand.Id = uint(brandId)
		shoe.Type = formData["type"].(string)
		shoe.Description = formData["description"].(string)
		shoe.Sku = formData["sku"].(string)
		shoe.Price = price
		shoe.Stock = stock
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
		temp, err := template.ParseFiles("views/shoes/index.html")
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
