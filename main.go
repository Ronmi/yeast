package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

func main() {
	var (
		data   string
		port   string
		ngconf string
	)
	flag.StringVar(&data, "data", "/var/lib/cheesecake/data.json", "path to store mapping")
	flag.StringVar(&port, "addr", ":8080", "address to listen")
	flag.StringVar(&ngconf, "conf", "/etc/nginx/sites-enabled/default", "path to nginx config")
	flag.Parse()

	p := NewPersistor(data, ngconf)
	if err := p.Load(); err != nil {
		log.Fatalf("Cannot load data from %s: %s", data, err)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Cannot load template: %s", err)
	}

	h := Handler{p, tmpl}
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/add", h.Add)
	http.HandleFunc("/del", h.Delete)
	http.HandleFunc("/enable", h.Enable)
	http.HandleFunc("/disable", h.Disable)

	log.Fatal(http.ListenAndServe(port, nil))
}
