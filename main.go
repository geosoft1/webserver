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
	})

	// setting a middleware (https://github.com/gorilla/mux#middleware)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	// https://www.google.com/webmasters/tools/home?hl=en

	// google site verification file
	// https://www.google.com/webmasters/verification/home?hl=en
	router.HandleFunc("/{google-site-verification}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.ServeFile(w, r, filepath.Join(path, vars["google-site-verification"]))
	})
	// google robots tester
	// https://www.google.com/webmasters/tools/robots-testing-tool
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(path, "robots.txt"))
	})

	http.ListenAndServe(":"+*port, router)
}
