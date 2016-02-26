package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// Handler handles all http requests
type Handler struct {
	Persistor *Persistor
	Template  *template.Template
}

func (h *Handler) err(w http.ResponseWriter, args ...interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, args...)
}

func (h *Handler) errf(w http.ResponseWriter, format string, args ...interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, format, args...)
}

// Index is default Handler
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.Template.Execute(w, h.Persistor.List())
}

// Add a mapping
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.err(w, "Cannot read input data!")
		return
	}

	name := r.FormValue("name")
	path := r.FormValue("path")
	upstream := r.FormValue("upstream")

	if upstream == "" || path == "" {
		h.err(w, "Path and Upstream cannot be empty string!")
		return
	}

	h.Persistor.Set(name, path, upstream)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Delete a mapping
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.err(w, "Cannot read input data!")
		return
	}

	name := r.FormValue("name")
	path := r.FormValue("path")

	if path == "" {
		h.err(w, "Path cannot be empty string!")
		return
	}

	h.Persistor.Unset(name, path)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Enable a mapping
func (h *Handler) Enable(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.err(w, "Cannot read input data!")
		return
	}

	name := r.FormValue("name")
	path := r.FormValue("path")

	if path == "" {
		h.err(w, "Path cannot be empty string!")
		return
	}

	h.Persistor.Enable(name, path)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Disable a mapping
func (h *Handler) Disable(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.err(w, "Cannot read input data!")
		return
	}

	name := r.FormValue("name")
	path := r.FormValue("path")

	if path == "" {
		h.err(w, "Path cannot be empty string!")
		return
	}

	h.Persistor.Disable(name, path)
	http.Redirect(w, r, "/", http.StatusFound)
}
