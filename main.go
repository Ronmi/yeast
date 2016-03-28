// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	var (
		data   string
		port   string
		ngconf string
		fend   string
		debug  bool
	)
	flag.StringVar(&data, "data", "/var/lib/cheesecake/data.json", "path to store mapping")
	flag.StringVar(&port, "addr", ":8080", "address to listen")
	flag.StringVar(&ngconf, "conf", "/etc/nginx/sites-enabled/default", "path to nginx config")
	flag.StringVar(&fend, "fe", "index.html", "Path to frontend file")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	p := NewPersistor(data, ngconf)
	if err := p.Load(); err != nil {
		log.Fatalf("Cannot load data from %s: %s", data, err)
	}

	tmpl, err := ioutil.ReadFile(fend)
	if err != nil {
		log.Fatalf("Cannot read frontend file from %s: %s", fend, err)
	}

	f := func() bool {
		cmd := exec.Command("nginx", "-s", "reload")
		return cmd.Run() == nil
	}

	if debug {
		f = func() bool {
			return true
		}
	}

	h := Handler{
		p,
		f,
	}
	http.HandleFunc("/api/list", h.List)
	http.HandleFunc("/api/create", h.Create)
	http.HandleFunc("/api/modify", h.Modify)
	http.HandleFunc("/api/delete", h.Delete)
	http.HandleFunc("/api/enable", h.Enable)
	http.HandleFunc("/api/disable", h.Disable)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(tmpl))
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
