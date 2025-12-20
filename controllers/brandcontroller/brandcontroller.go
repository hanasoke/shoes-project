package brandcontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
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
	t = template.Must(t.ParseFiles("views/brand/index.html"))

	t.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("views/brand/create.html"))
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
			temp := template.Must(template.ParseFiles("views/brand/create.html"))
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
			temp := template.Must(template.ParseFiles("views/brand/create.html"))
			temp.Execute(w, data)
			return
		}

		http.Redirect(w, r, "/brands", http.StatusSeeOther)
	}
}
