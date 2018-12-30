package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var (
	port = flag.String("port", "8080", "port")
)

var templ = template.New("templ")

func main() {
	flag.Parse()
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	templ.ParseGlob(filepath.Join(path, "www", "*.html"))
	router := mux.NewRouter()
	router.PathPrefix("/www/").Handler(http.StripPrefix("/www/", http.FileServer(http.Dir(filepath.Join(path, "www")))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "index", nil)
		log.Println(r.RemoteAddr, r.RequestURI)
	})

	http.ListenAndServe(":"+*port, router)
}
