package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/models/shoemodel"
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
}
