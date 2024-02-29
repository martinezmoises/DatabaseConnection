package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to quotebox."))
}

func createQuoteForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("/home/moises/Downloads/Last Class/quote_form_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/quote", createQuoteForm)
	log.Println("Starting server on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
