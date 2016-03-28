// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import (
	"encoding/json"
	"net/http"
)

// Handler handles all api calls
type Handler struct {
	Persistor   *Persistor
	ReloadNginx func() bool // reload nginx, return true if success
}

// List lists all known mapping data
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	data := h.Persistor.List()
	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}

	w.Write(buf)
}

// Set add or overwrite a mapping data
func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")
	upstream := r.PostFormValue("upstream")
	custom := r.PostFormValue("custom_tags")

	if name == "" || path == "" || upstream == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you must pass at least name, path and upstream"))
		return
	}

	data := h.Persistor.Set(name, path, upstream, custom)
	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}

	if !h.ReloadNginx() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot reload Nginx."))
		return
	}

	w.Write(buf)
}

// Unset removes a existing mapping
func (h *Handler) Unset(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")

	if name == "" || path == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you must pass at least name and path"))
		return
	}

	if !h.Persistor.Unset(name, path) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such data"))
		return
	}

	if !h.ReloadNginx() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot reload Nginx."))
		return
	}

}

// Enable enables some of known data
func (h *Handler) Enable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")

	data := h.Persistor.Enable(name, path)

	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}

	if !h.ReloadNginx() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot reload Nginx."))
		return
	}

	w.Write(buf)
}

// Disable disables some of known data
func (h *Handler) Disable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")

	data := h.Persistor.Disable(name, path)

	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}

	if !h.ReloadNginx() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot reload Nginx."))
		return
	}

	w.Write(buf)
}
