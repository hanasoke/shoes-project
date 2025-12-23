package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/shoemodel"
	"strconv"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	shoes := shoemodel.GetAll
	data := map[string]any{
		"shoes": shoes,
	}

	temp, err := template.ParseFiles("views/shoes/index.html")

	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		temp, err := template.ParseFiles("views/shoes/create.html")
		if err != nil {
			panic(err)
		}

		temp.Execute(w, nil)
	}

	if r.Method == "POST" {
		var shoe entities.Shoe

		brandId, err := strconv.Atoi(r.FormValue("brand_id"))
		if err != nil {
			panic(err)
		}

		stock, err := strconv.Atoi(r.FormValue("shoe_stock"))
		if err != nil {
			panic(err)
		}

		shoe.Name = r.FormValue("shoe_name")
		shoe.Brand.Brand_Id = uint(brandId)
		shoe.Stock = int64(stock)
		shoe.Description = r.FormValue("shoe_description")
		shoe.CreatedAt = time.Now()

		if ok := shoemodel.Create(shoe); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
