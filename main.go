package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	port = flag.String("p", "8080", "port")
)

func main() {
	http.HandleFunc("/", home)
	println("started on http://localhost:" + *port)
	e := http.ListenAndServe(":"+*port, nil)
	if e != nil {
		log.Println(e)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, filename)
	case http.MethodPut:
		os.MkdirAll(filepath.Dir(filename), 0755)
		fo, e := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if e != nil {
			http.Error(w, e.Error(), 500)
			log.Println(e)
			return
		}
		defer fo.Close()
		defer r.Body.Close()
		_, e = io.Copy(fo, r.Body)
		if e != nil {
			log.Println(e)
			http.Error(w, e.Error(), 500)
			return
		}
		w.Write([]byte("OK"))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
