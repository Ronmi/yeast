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
	res := h.Persistor.List()
	data := map[string]map[string]*Mapping{}
	for k, v := range res {
		data[k] = v.Paths
	}
	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}

	w.Write(buf)
}

// Create add a mapping data
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	res := h.Persistor.Create(name, path, upstream, custom)
	if res == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	data := map[string]map[string]*Mapping{
		res.ServerName: res.Paths,
	}
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

// Modify a mapping data
func (h *Handler) Modify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")
	newPath := r.PostFormValue("new_path")
	upstream := r.PostFormValue("new_upstream")
	custom := r.PostFormValue("new_custom_tags")

	if name == "" || path == "" || newPath == "" || upstream == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you must pass at least name, path, new_path and new_upstream"))
		return
	}

	res := h.Persistor.Modify(name, path, newPath, upstream, custom)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := map[string]map[string]*Mapping{
		res.ServerName: res.Paths,
	}
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

// Delete a existing mapping
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")

	if name == "" || path == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you must pass at least name and path"))
		return
	}

	res := h.Persistor.Delete(name, path)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such data"))
		return
	}

	if !h.ReloadNginx() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot reload Nginx."))
		return
	}

	data := map[string]map[string]*Mapping{
		res.ServerName: res.Paths,
	}
	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot serialize data to json format."))
		return
	}
	w.Write(buf)
}

// Enable enables some of known data
func (h *Handler) Enable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.PostFormValue("name")
	path := r.PostFormValue("path")

	res := h.Persistor.Enable(name, path)

	data := map[string]map[string]*Mapping{}
	for k, v := range res {
		data[k] = v.Paths
	}
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

	res := h.Persistor.Disable(name, path)

	data := map[string]map[string]*Mapping{}
	for k, v := range res {
		data[k] = v.Paths
	}
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
