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

		shoe.Name = r.FormValue("name")
		shoe.Brand.Id = uint(brandId)
		shoe.Type = r.FormValue("type")
		shoe.Description = r.FormValue("description")
		shoe.SKU = r.FormValue("sku")
		shoe.Price = int64(price)
		shoe.Stock = int64(stock)
		shoe.CreatedAt = time.Now()

		if ok := shoemodel.Create(shoe); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
