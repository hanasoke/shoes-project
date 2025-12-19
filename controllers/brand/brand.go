package brand

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brand"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	brands := brand.GetAll()
	data := map[string]any{
		"brands": brands,
	}

	temp, err := template.ParseFiles("views/brand/index.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/brand/create.html")
		if err != nil {
			panic(err)
		}

		temp.Execute(w, nil)
	}

	if r.Method == "POST" {
		var brand entities.Brand

		brand.Brand_Name = r.FormValue("brand_name")
		brand.CreatedAt = time.Now()
		brand.UpdatedAt = time.Now()

		// if ok := brand.Created(brand); !ok {
		// 	temp, _ := template.ParseFiles("views/brand/create.html")
		// 	temp.Execute(w, nil)
		// }

		http.Redirect(w, r, "/brands", http.StatusSeeOther)

	}
}
