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

// Fungsi untuk format Rupiah
func formatRupiah(amount int64) string {
	str := strconv.FormatInt(amount, 10)

	// Format dengan pemisah ribuan
	n := len(str)
	if n <= 3 {
		return "Rp " + str
	}

	var result []string
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}

	return "Rp " + strings.Join(result, ".")
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Buat funcMap dengan fungsi helper
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatRupiah": formatRupiah,
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
	// Ambil ID dari query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Konversi ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Ambil data shoe dari model
	shoe, err := shoemodel.Detail(id)
	if err != nil {
		// Handle jika shoe tidak ditemukan
		if err.Error() == "shoe not found" {
			http.Error(w, "Shoe not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Buat funcMap untuk detail
	funcMap := template.FuncMap{
		"formatRupiah": formatRupiah,
	}

	// Parse template dengan fungsi helper
	t := template.New("detail.html").Funcs(funcMap)
	t, err = t.ParseFiles("views/shoes/detail.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Prepare data untuk template
	data := map[string]interface{}{
		"shoe": shoe,
	}

	// Execute template dengan data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
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
			"errors": map[string]string{
				"name":        "",
				"brand":       "",
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

		// Inisialisasi map untuk error per field
		fieldErrors := map[string]string{
			"name":        "",
			"brand":       "",
			"type":        "",
			"description": "",
			"sku":         "",
			"price":       "",
			"stock":       "",
		}

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
			fieldErrors["name"] = "Shoe name cannot be empty"
		}

		// Validasi brand
		hasBrandError := false
		if idBrandStr == "" || idBrandStr == "#" {
			fieldErrors["brand"] = "Brand is required"
			hasBrandError = true
			idBrandStr = "" // Set kosong untuk konsistensi
		}

		// Validasi type
		if shoeType == "" {
			fieldErrors["type"] = "Type cannot be empty"
		}

		// Validasi description
		if description == "" {
			fieldErrors["description"] = "Description cannot be empty"
		}

		// Validasi sku
		if sku == "" {
			fieldErrors["sku"] = "SKU cannot be empty"
		}

		// Validasi price
		if priceStr == "" {
			fieldErrors["price"] = "Price cannot be empty"
		}

		// Validasi stock
		if stockStr == "" {
			fieldErrors["stock"] = "Stock cannot be empty"
		}

		// Konversi dan validasi tipe data
		var idBrand int
		var price, stock int64
		var hasValidationError bool

		// Cek apakah ada error validasi field required
		for _, errMsg := range fieldErrors {
			if errMsg != "" {
				hasValidationError = true
				break
			}
		}

		if !hasBrandError && idBrandStr != "" && idBrandStr != "#" {
			id, err := strconv.Atoi(idBrandStr)
			if err != nil {
				fieldErrors["brand"] = "Invalid brand ID"
				hasValidationError = true
			} else if id <= 0 {
				fieldErrors["brand"] = "Brand ID must be positive"
				hasValidationError = true
			} else {
				idBrand = id
				formData["idBrand"] = uint(id)
			}
		}

		// Validasi format price jika tidak kosong
		if priceStr != "" && fieldErrors["price"] == "" {
			p, err := strconv.ParseInt(priceStr, 10, 64)
			if err != nil {
				fieldErrors["price"] = "Price must be a valid number"
				hasValidationError = true
			} else if p <= 0 {
				fieldErrors["price"] = "Price must be a positive number"
				hasValidationError = true
			} else {
				price = p
			}
		}

		// Validasi format stock jika tidak kosong
		if stockStr != "" && fieldErrors["stock"] == "" {
			s, err := strconv.ParseInt(stockStr, 10, 64)
			if err != nil {
				fieldErrors["stock"] = "Stock must be a valid number"
				hasValidationError = true
			} else if s < 0 {
				fieldErrors["stock"] = "Stock must be zero or positive number"
				hasValidationError = true
			} else {
				stock = s
			}
		}

		// Jika ada error, tampilkan form kembali
		if hasValidationError {
			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"errors": fieldErrors,
				"form":   formData,
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
			UpdatedAt:   time.Now(),
		}

		// Simpan ke database
		err := shoemodel.Create(shoe)
		if err != nil {
			msg := "Failed to create shoe"
			if err == shoemodel.ErrDuplicateShoe {
				msg = "Shoe already exists"
			}

			// Untuk error duplikat, tampilkan di field name
			fieldErrors["name"] = msg

			brands := brandmodel.GetAll()
			data := map[string]any{
				"brands": brands,
				"errors": fieldErrors,
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

func Edit(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Konversi ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		// Ambil data shoe dari model
		shoe, err := shoemodel.Detail(id)
		if err != nil {
			// Handle jika shoe tidak ditemukan
			if err.Error() == "shoe not found" {
				http.Error(w, "Shoe not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Ambil semua brand untuk dropdown
		brands := brandmodel.GetAll()

		// Buat funcMap
		funcMap := template.FuncMap{
			"formatRupiah": formatRupiah,
		}

		// Parse template
		t := template.New("edit.html").Funcs(funcMap)
		t, err = t.ParseFiles("views/shoes/edit.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		// Prepare data untuk template
		data := map[string]interface{}{
			"shoe":   shoe,
			"brands": brands,
			"form": map[string]interface{}{
				"name":        shoe.Name,
				"idBrand":     shoe.Brand.Id,
				"type":        shoe.Type,
				"description": shoe.Description,
				"sku":         shoe.SKU,
				"price":       shoe.Price,
				"stock":       shoe.Stock,
			},
			"errors": map[string]string{
				"name":        "",
				"brand":       "",
				"type":        "",
				"description": "",
				"sku":         "",
				"price":       "",
				"stock":       "",
			},
		}

		// Execute template dengan data
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
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

		// Inisialisasi map untuk error per field
		fieldErrors := map[string]string{
			"name":        "",
			"brand":       "",
			"type":        "",
			"description": "",
			"sku":         "",
			"price":       "",
			"stock":       "",
		}

		// Data form untuk dikembalikan ke template
		formData := map[string]interface{}{
			"name":        name,
			"type":        shoeType,
			"description": description,
			"sku":         sku,
			"price":       priceStr,
			"stock":       stockStr,
		}

		// Validasi required fields
		hasValidationError := false

		// Validasi nama
		if name == "" {
			fieldErrors["name"] = "Shoe name cannot be empty"
			hasValidationError = true
		}

		// Validasi brand
		if idBrandStr == "" || idBrandStr == "#" {
			fieldErrors["brand"] = "Brand is required"
			hasValidationError = true
		}

		// Validasi type
		if shoeType == "" {
			fieldErrors["type"] = "Type cannot be empty"
			hasValidationError = true
		}

		// Validasi description
		if description == "" {
			fieldErrors["description"] = "Description cannot be empty"
			hasValidationError = true
		}

		// Validasi sku
		if sku == "" {
			fieldErrors["sku"] = "SKU cannot be empty"
			hasValidationError = true
		}

		// Validasi price
		if priceStr == "" {
			fieldErrors["price"] = "Price cannot be empty"
			hasValidationError = true
		}

		// Validasi stock
		if stockStr == "" {
			fieldErrors["stock"] = "Stock cannot be empty"
			hasValidationError = true
		}

		// Konversi dan validasi tipe data
		var idBrand uint

		if idBrandStr != "" && idBrandStr != "#" && fieldErrors["brand"] == "" {
			id, err := strconv.Atoi(idBrandStr)
			if err != nil {
				fieldErrors["brand"] = "Invalid brand selection"
				hasValidationError = true
			} else if id <= 0 {
				fieldErrors["brand"] = "Please select a valid brand"
				hasValidationError = true
			} else {
				idBrand = uint(id)
				formData["idBrand"] = idBrand
			}
		}

		// Validasi price - dengan pesan error yang lebih spesifik
		var price int64
		if priceStr == "" {
			fieldErrors["price"] = "Price cannot be empty"
			hasValidationError = true
		} else {
			p, err := strconv.ParseInt(priceStr, 10, 64)
			if err != nil {
				// Deteksi tipe error
				if strings.Contains(err.Error(), "value out of range") {
					fieldErrors["price"] = "Price value is too large"
				} else if strings.Contains(err.Error(), "invalid syntax") {
					fieldErrors["price"] = "Price must be a valid number (e.g., 150000)"
				} else {
					fieldErrors["price"] = "Invalid price format"
				}
				hasValidationError = true
			} else if p <= 0 {
				fieldErrors["price"] = "Price must be greater than 0"
				hasValidationError = true
			} else if p > 1000000000 { // Batas maksimal 1 miliar
				fieldErrors["price"] = "Price cannot exceed Rp 1,000,000,000"
				hasValidationError = true
			} else {
				price = p
				formData["price"] = price // Simpan sebagai integer, bukan string
			}
		}

		// Validasi stock - dengan pesan error yang lebih spesifik
		var stock int64
		if stockStr == "" {
			fieldErrors["stock"] = "Stock cannot be empty"
			hasValidationError = true
		} else {
			s, err := strconv.ParseInt(stockStr, 10, 64)
			if err != nil {
				if strings.Contains(err.Error(), "value out of range") {
					fieldErrors["stock"] = "Stock value is too large"
				} else if strings.Contains(err.Error(), "invalid syntax") {
					fieldErrors["stock"] = "Stock must be a valid whole number"
				} else {
					fieldErrors["stock"] = "Invalid stock format"
				}
				hasValidationError = true
			} else if s < 0 {
				fieldErrors["stock"] = "Stock cannot be negative"
				hasValidationError = true
			} else if s > 100000 { // Batas maksimal stok
				fieldErrors["stock"] = "Stock cannot exceed 100,000 units"
				hasValidationError = true
			} else {
				stock = s
				formData["stock"] = stock // Simpan sebagai integer, bukan string
			}
		}

		// Jika ada error validasi, tampilkan form kembali
		if hasValidationError {
			brands := brandmodel.GetAll()
			shoe, _ := shoemodel.Detail(id)

			data := map[string]interface{}{
				"shoe":   shoe,
				"brands": brands,
				"errors": fieldErrors,
				"form":   formData,
			}

			funcMap := template.FuncMap{
				"formatRupiah": formatRupiah,
			}

			t := template.New("edit.html").Funcs(funcMap)
			t, err = t.ParseFiles("views/shoes/edit.html")
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

			t.Execute(w, data)
			return
		}

		// Buat objek shoe untuk update
		shoeUpdate := entities.ShoeUpdate{
			Name:        name,
			IdBrand:     idBrand,
			Type:        shoeType,
			Description: description,
			SKU:         sku,
			Price:       price,
			Stock:       stock,
			UpdatedAt:   time.Now(),
		}

		// Update data ke database
		err = shoemodel.Update(id, shoeUpdate)
		if err != nil {
			msg := "Failed to update shoe"
			if err == shoemodel.ErrDuplicateShoe {
				msg = "Shoe name already exists"
			}

			// Tampilkan form kembali dengan error
			brands := brandmodel.GetAll()
			shoe, _ := shoemodel.Detail(id)

			fieldErrors["name"] = msg
			data := map[string]interface{}{
				"shoe":   shoe,
				"brands": brands,
				"errors": fieldErrors,
				"form":   formData,
			}

			funcMap := template.FuncMap{
				"formatRupiah": formatRupiah,
			}

			t := template.New("edit.html").Funcs(funcMap)
			t, err = t.ParseFiles("views/shoes/edit.html")
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

			t.Execute(w, data)
			return
		}

		// Redirect ke halaman shoes dengan pesan sukses
		http.Redirect(w, r, "/shoes?success=updated", http.StatusSeeOther)
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
