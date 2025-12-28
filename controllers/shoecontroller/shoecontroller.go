package shoecontroller

import (
	"html/template"
	"net/http"
	"shoes-project/entities"
	"shoes-project/models/brandmodel"
	"shoes-project/models/shoemodel"
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

// Struct to pass data to template
type FormData struct {
	Brands  []entities.Brand
	Errors  []shoemodel.ValidationError
	OldData map[string]string
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/shoes/create.html")
		if err != nil {
			panic(err)
		}

		brands := brandmodel.GetAll()
		data := FormData{
			Brands:  brands,
			Errors:  nil,
			OldData: make(map[string]string),
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var shoe entities.Shoe

		// Get form values and store for repopulation
		oldData := make(map[string]string)
		oldData["name"] = r.FormValue("name")
		oldData["type"] = r.FormValue("type")

		http.Redirect(w, r, "/shoes", http.StatusSeeOther)
	}
}
