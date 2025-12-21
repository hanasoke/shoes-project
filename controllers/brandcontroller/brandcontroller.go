package brandcontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"strconv"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	brands := brandmodel.GetAll()
	data := map[string]any{
		"brands": brands,
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
			UpdatedAt:  time.Now(),
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

		http.Redirect(w, r, "/brands", http.StatusSeeOther)
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
		var brand entities.Brand

		idString := r.FormValue("brand_id")
		brand_id, err := strconv.Atoi(idString)
		if err != nil {
			panic(err)
		}

		brand.Brand_Name = r.FormValue("brand_name")
		brand.UpdatedAt = time.Now()

		if ok := brandmodel.Update(brand_id, brand); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/brands", http.StatusSeeOther)
	}
}
