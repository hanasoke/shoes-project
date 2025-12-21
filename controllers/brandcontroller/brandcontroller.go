package brandcontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
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

	brands := brandmodel.GetAll()
	data := map[string]any{
		"brands":  brands,
		"success": success,
	}

	t := template.New("index.html").Funcs(funcMap)
	t = template.Must(t.ParseFiles("views/brands/index.html"))

	t.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("views/brands/create.html"))
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		brandName := r.FormValue("brand_name")

		// 1️⃣ NULL / empty validation
		if brandName == "" {
			data := map[string]any{
				"error": "Brand name cannot be empty",
			}
			temp := template.Must(template.ParseFiles("views/brands/create.html"))
			temp.Execute(w, data)
			return
		}

		brand := entities.Brand{
			Brand_Name: brandName,
			CreatedAt:  time.Now(),
		}

		err := brandmodel.Create(brand)
		if err != nil {
			msg := "Failed to create brand"

			if err == brandmodel.ErrDuplicateBrand {
				msg = "Brand already exists"
			}

			data := map[string]any{
				"error": msg,
			}
			temp := template.Must(template.ParseFiles("views/brands/create.html"))
			temp.Execute(w, data)
			return
		}

		http.Redirect(w, r, "/brands?success=created", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/brands/edit.html")
		if err != nil {
			panic(err)
		}

		idString := r.URL.Query().Get("brand_id")
		brand_id, err := strconv.Atoi(idString)
		if err != nil {
			panic(err)
		}

		brand := brandmodel.Detail(brand_id)
		data := map[string]any{
			"brand": brand,
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		idString := r.FormValue("brand_id")
		brand_id, _ := strconv.Atoi(idString)

		brandName := strings.TrimSpace(r.FormValue("brand_name"))

		// ❌ NULL validation
		if brandName == "" {
			brand := brandmodel.Detail(brand_id)
			data := map[string]any{
				"brand": brand,
				"error": "Brand name cannot be empty",
			}

			temp := template.Must(template.ParseFiles("views/brands/edit.html"))
			temp.Execute(w, data)
			return
		}

		brand := entities.Brand{
			Brand_Name: brandName,
			UpdatedAt:  time.Now(),
		}

		err := brandmodel.Update(brand_id, brand)
		if err != nil {
			msg := "Failed to update brand"

			if err == brandmodel.ErrDuplicateBrand {
				msg = "Brand name already exists"
			}

			oldBrand := brandmodel.Detail(brand_id)
			data := map[string]any{
				"brand": oldBrand,
				"error": msg,
			}

			temp := template.Must(template.ParseFiles("views/brands/edit.html"))
			temp.Execute(w, data)
			return
		}

		http.Redirect(w, r, "/brands?success=updated", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("brand_id")
	brand_id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	if err := brandmodel.Delete(brand_id); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/brands?success=deleted", http.StatusSeeOther)
}
