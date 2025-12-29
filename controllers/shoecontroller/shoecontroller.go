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

func Detail(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	shoe := shoemodel.Detail(id)
	data := map[string]any{
		"shoe": shoe,
	}

	temp, err := template.ParseFiles("views/shoes/detail.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		brands := brandmodel.GetAll()
		data := map[string]any{
			"brands": brands,
		}

		temp := template.Must(template.ParseFiles("views/shoes/create.html"))

		temp.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		// Ambil data dari form
		name := strings.TrimSpace(r.FormValue("name"))
		idBrandStr := r.FormValue("id_brand")
		shoeType := strings.TrimSpace(r.FormValue("type"))
		description := strings.TrimSpace(r.FormValue("description"))
		sku := strings.TrimSpace(r.FormValue("sku"))
		priceStr := r.FormValue("price")
		stockStr := r.FormValue("stock")

		// Validasi required fields
		var validationErrors []string

		if name == "" {
			validationErrors = append(validationErrors, "Shoe name cannot be empty")
		}

		if idBrandStr == "" {
			validationErrors = append(validationErrors, "Brand is required")
		}

		if shoeType == "" {
			validationErrors = append(validationErrors, "Type cannot be empty")
		}

		if description == "" {
			validationErrors = append(validationErrors, "Type cannot be empty")
		}

		if sku == "" {
			validationErrors = append(validationErrors, "SKU cannot be empty")
		}

		if priceStr == "" {
			validationErrors = append(validationErrors, "Price cannot be empty")
		}

		if stockStr == "" {
			validationErrors = append(validationErrors, "Price cannot be empty")
		}

		if len(validationErrors) > 0 {
			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"error":  strings.Join(validationErrors, ","),
			}
			temp := template.Must(template.ParseFiles("views/shoes/create.html"))
			temp.Execute(w, data)
			return
		}

		// Konversi tipe data
		idBrand, err := strconv.Atoi(idBrandStr)
		if err != nil {
			validationErrors = append(validationErrors, "Invalid brand ID")
		}

		price, err := strconv.ParseInt(priceStr, 10, 64)
		if err != nil || price <= 0 {
			validationErrors = append(validationErrors, "Price must be a positive number")
		}

		stock, err := strconv.ParseInt(stockStr, 10, 64)
		if err != nil || stock < 0 {
			validationErrors = append(validationErrors, "Stock must be zero or positive number")
		}

		if len(validationErrors) > 0 {
			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"error":  strings.Join(validationErrors, ","),
			}
			temp := template.Must(template.ParseFiles("views/shoes/create.html"))
			temp.Execute(w, data)
			return
		}

		// Buat objek shoe
		shoe := entities.ShoeCreate{
			Name:        name,
			IdBrand:     uint(idBrand),
			Type:        shoeType,
			Description: description,
			SKU:         sku,
			Price:       price,
			Stock:       stock,
			CreatedAt:   time.Now(),
		}

		// Simpan ke database
		err = shoemodel.Create(shoe)
		if err != nil {
			msg := "Failed to create shoe"
			if err == shoemodel.ErrDuplicateShoe {
				msg = "Shoe already exists"
			}

			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"error":  msg,
			}
			temp := template.Must(template.ParseFiles("views/shoes/create.html"))

			temp.Execute(w, data)
			return
		}

		http.Redirect(w, r, "/shoes?success=created", http.StatusSeeOther)
	}
}
