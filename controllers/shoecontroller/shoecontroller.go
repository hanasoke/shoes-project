package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"shoes-project/models/shoemodel"
	"strconv"
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
		temp := template.Must(template.ParseFiles("views/shoes/create.html"))

		brands := brandmodel.GetAll()
		data := map[string]any{
			"brands": brands,
		}

		temp.Execute(w, data)
	}

	if r.Method == http.MethodPost {

		brandId, err := strconv.Atoi(r.FormValue("id_brand"))
		if err != nil {
			panic(err)
		}

		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			panic(err)
		}

		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			panic(err)
		}

		shoeName := r.FormValue("name")
		BrandId := uint(brandId)
		Type = r.FormValue("type")
		Description = r.FormValue("description")
		SKU := r.FormValue("sku")

		// 1️⃣ NULL / empty validation
		if shoeName == "" {
			data := map[string]any{
				"error": "Shone Name cannot be empty",
			}
			temp := template.Must(template.ParseFiles("views/shoes/create.html"))
			temp.Execute(w, data)
			return
		}

		brand := entities.Shoe{
			Name:      shoeName,
			Brand:     BrandId,
			CreatedAt: time.Now(),
			Price:     int64(price),
			Stock:     int64(stock),
		}

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
