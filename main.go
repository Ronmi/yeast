// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/Patrolavia/toolkit/session"
)

func main() {
	var (
		data   string
		port   string
		ngconf string
		fend   string
		pass   string
		debug  bool
	)
	flag.StringVar(&data, "data", "/var/lib/cheesecake/data.json", "path to store mapping")
	flag.StringVar(&port, "addr", ":8080", "address to listen")
	flag.StringVar(&ngconf, "conf", "/etc/nginx/sites-enabled/default", "path to nginx config")
	flag.StringVar(&fend, "fe", ".", "Path to directory holding frontend files")
	flag.StringVar(&pass, "pass", "", "password to lock the manage page")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	p := NewPersistor(data, ngconf)
	if err := p.Load(); err != nil {
		log.Fatalf("Cannot load data from %s: %s", data, err)
	}

	tmpl, err := ioutil.ReadFile(fend + "/index.html")
	if err != nil {
		log.Fatalf("Cannot read index page from %s/index.html: %s", fend, err)
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

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(tmpl)
	}
	realRoot := rootHandler
	loginPage, err := ioutil.ReadFile(fend + "/login.html")
	if err != nil {
		log.Fatalf("Cannot read login page from %s/login.html: %s", fend, err)
	}

	if pass != "" {
		realRoot = (&session.Middleware{
			Manager: &session.Manager{},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				sess := r.Context().Value("session").(*session.Session)
				if sess == nil {
					// no session, show login page
					w.Write(loginPage)
					return
				}
				defer sess.Save(w, session.DefaultCookieMaker)

				if sess.Data() == "ok" {
					rootHandler(w, r)
					return
				}
				r.ParseForm()
				if r.PostFormValue("pass") != pass {
					// incorrect password, show login page
					w.Write(loginPage)
					return
				}

				sess.SetData("ok")
				rootHandler(w, r)
			},
		}).Handle
	}
	http.HandleFunc("/", realRoot)

	log.Fatal(http.ListenAndServe(port, nil))
}
