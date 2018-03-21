package shortener

import (
	"net/http"
	"fmt"
	"html/template"
)

// Handler for <host>/
// Returns content of index.html
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/index.html")
}

// Handler for <host>/create
// Returns content of create.html
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	v := r.Form.Get("longurl")

	tmpl := "tmpl/create.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Printf("Could not parse template file '%v'. %v\n", tmpl, err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	s, err := generateShortUrl(v)
	if err != nil {
		fmt.Printf("Could not generate short url. %v\n", err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, s)
	if err != nil {
		fmt.Printf("Could not execute template. %v\n", err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Path[5:]

	longUrl, err := getLongUrl(alias)
	if err != nil {
		fmt.Printf("Could not get long url. %v\n", err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusSeeOther)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	http.ServeFile(w, r, "tmpl/404.html")
}
