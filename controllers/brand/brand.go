package brand

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/brand/index.html")

	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/brand/create.html")
		if err != nil {
			panic(err)
		}

		temp.Execute(w, nil)
	}
}
