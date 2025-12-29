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

	temp, err := template.ParseFiles("views/shoes/detail.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		brands := brandmodel.GetAll()
		data := map[string]any{
			"brands": brands,
			"form": map[string]any{
				"name":        "",
				"idBrand":     0,
				"type":        "",
				"description": "",
				"sku":         "",
				"price":       "",
				"stock":       "",
			},
		}

		temp := template.Must(template.ParseFiles("views/shoes/create.html"))

		temp.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		// Ambil data dari form
		name := strings.TrimSpace(r.FormValue("name"))
		idBrandStr := strings.TrimSpace(r.FormValue("id_brand"))
		shoeType := strings.TrimSpace(r.FormValue("type"))
		description := strings.TrimSpace(r.FormValue("description"))
		sku := strings.TrimSpace(r.FormValue("sku"))
		priceStr := strings.TrimSpace(r.FormValue("price"))
		stockStr := strings.TrimSpace(r.FormValue("stock"))

		// Validasi required fields
		var validationErrors []string
		var errorBrand string

		// Data form untuk dikembalikan ke template
		formData := map[string]any{
			"name":        name,
			"type":        shoeType,
			"description": description,
			"sku":         sku,
			"price":       priceStr,
			"stock":       stockStr,
		}

		// Validasi nama
		if name == "" {
			validationErrors = append(validationErrors, "Shoe name cannot be empty")
		}

		// Validasi brand - ini adalah masalah utama
		if idBrandStr == "" {
			validationErrors = append(validationErrors, "Brand is required")
			errorBrand = "Brand is required"
		} else if idBrandStr == "#" {
			// Periksa jika user memilih placeholder
			validationErrors = append(validationErrors, "Brand is required")
			errorBrand = "Brand is required"
			idBrandStr = "" // Set kosong untuk konsistensi

		}

		// Validasi field lainnya
		if shoeType == "" {
			validationErrors = append(validationErrors, "Type cannot be empty")
		}

		if description == "" {
			validationErrors = append(validationErrors, "Description cannot be empty")
		}

		if sku == "" {
			validationErrors = append(validationErrors, "SKU cannot be empty")
		}

		if priceStr == "" {
			validationErrors = append(validationErrors, "Price cannot be empty")
		}

		if stockStr == "" {
			validationErrors = append(validationErrors, "Stock cannot be empty")
		}

		// Konversi dan validasi tipe data
		var idBrand int
		var price, stock int64
		var conversionErrors []string

		if idBrandStr != "" && idBrandStr != "#" {
			id, err := strconv.Atoi(idBrandStr)
			if err != nil {
				conversionErrors = append(conversionErrors, "Invalid brand ID")
				errorBrand = "Invalid brand ID"
			} else if id <= 0 {
				conversionErrors = append(conversionErrors, "Brand ID must be positive")
				errorBrand = "Brand ID must be positive"
			} else {
				idBrand = id
				formData["idBrand"] = uint(id)
			}
		}

		if priceStr != "" {
			p, err := strconv.ParseInt(priceStr, 10, 64)
			if err != nil {
				conversionErrors = append(conversionErrors, "Price must be a valid number")
			} else if p <= 0 {
				conversionErrors = append(conversionErrors, "Price must be a positive number")
			} else {
				price = p
			}
		}

		if stockStr != "" {
			s, err := strconv.ParseInt(stockStr, 10, 64)
			if err != nil {
				conversionErrors = append(conversionErrors, "Stock must be a valid number")
			} else if s < 0 {
				conversionErrors = append(conversionErrors, "Stock must be zero or positive number")
			} else {
				stock = s
			}
		}

		// Gabungkan semua error
		allErrors := append(validationErrors, conversionErrors...)

		if len(allErrors) > 0 {
			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands":     brands,
				"error":      strings.Join(allErrors, ", "),
				"errorBrand": errorBrand,
				"form":       formData,
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
		err := shoemodel.Create(shoe)
		if err != nil {
			msg := "Failed to create shoe"
			if err == shoemodel.ErrDuplicateShoe {
				msg = "Shoe already exists"
			}

			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"error":  msg,
				"form":   formData,
			}
			temp := template.Must(template.ParseFiles("views/shoes/create.html"))

			temp.Execute(w, data)
			return
		}

		// Redirect ke halaman index dengan pesan sukses
		http.Redirect(w, r, "/shoes?success=created", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	if err := shoemodel.Delete(id); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/shoes?success=deleted", http.StatusSeeOther)
}
