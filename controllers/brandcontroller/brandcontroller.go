package brandcontroller

import (
	"html/template"
	"net/http"
	"shoes-project/models/brandmodel"
)

func Index(w http.ResponseWriter, r *http.Request) {
	brands := brandmodel.GetAll()
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
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("views/brand/create.html"))
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {

	}
}
