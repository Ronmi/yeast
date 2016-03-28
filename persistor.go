// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

// Persistor holds all server info and save/load it into disk
type Persistor struct {
	filename string
	conffile string
	servers  map[string]*NginxServer
	*sync.Mutex
}

// NewPersistor creates a persistor
func NewPersistor(fn, conf string) *Persistor {
	return &Persistor{
		fn,
		conf,
		map[string]*NginxServer{},
		&sync.Mutex{},
	}
}

// Save configs to file
func (p *Persistor) Save() (err error) {
	p.Lock()
	defer p.Unlock()

	return p.doSave()
}

func (p *Persistor) doSave() (err error) {
	buf := make([]*NginxServer, 0, len(p.servers))
	for _, svr := range p.servers {
		buf = append(buf, svr)
	}

	str, err := json.Marshal(buf)
	if err != nil {
		return
	}

	f, err := os.Create(p.filename)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.Write(str); err != nil {
		return
	}

	if err = p.export(); err != nil {
		return
	}

	cmd := exec.Command("nginx", "-s", "reload")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return
}

func (p *Persistor) export() error {
	f, err := os.Create(p.conffile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, srv := range p.servers {
		fmt.Fprintln(f, srv.Export())
	}
	return nil
}

// Load configs from file
func (p *Persistor) Load() (err error) {
	p.Lock()
	defer p.Unlock()

	data, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return
	}

	var buf []NginxServer
	err = json.Unmarshal(data, &buf)
	if err != nil {
		return
	}

	p.servers = map[string]*NginxServer{}
	for _, srv := range buf {
		srv := srv
		for _, mapping := range srv.Paths {
			mapping.Enabled = false
		}
		p.servers[srv.ServerName] = &srv
	}

	return
}

func (p *Persistor) getServer(name string) *NginxServer {
	ret, ok := p.servers[name]
	if !ok {
		ret = NewServer(name)
		p.servers[name] = ret
	}

	return ret
}

// Create a path to upstream mapping
func (p *Persistor) Create(name, path, upstream, custom string) (ret *NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	srv := p.getServer(name)
	if srv.Create(path, upstream, custom) {
		ret = srv
	}

	return
}

// Modify a path to upstream mapping
func (p *Persistor) Modify(name, path, newPath, upstream, custom string) (ret *NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	srv := p.getServer(name)
	if srv.Modify(path, newPath, upstream, custom) {
		ret = srv
	}

	return
}

// Delete a path-upstream mapping
func (p *Persistor) Delete(name, path string) (ret *NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	if _, ok := p.servers[name]; !ok {
		return
	}

	ret = p.getServer(name)
	if !ret.Delete(path) {
		return
	}

	if ret.Len() < 1 {
		delete(p.servers, name)
	}
	return
}

// Enable a mapping
func (p *Persistor) Enable(name, path string) (ret map[string]*NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	if name == "" {
		return p.enableAll()
	}

	var srv *NginxServer
	if path == "" {
		srv = p.enableServer(name)
	} else {
		srv = p.enableOne(name, path)
	}

	ret = map[string]*NginxServer{
		srv.ServerName: srv,
	}
	return
}

func (p *Persistor) enableOne(name, path string) (ret *NginxServer) {
	ret = p.getServer(name)
	ret.Enable(path)
	return
}

func (p *Persistor) enableServer(name string) (ret *NginxServer) {
	ret = p.getServer(name)
	for path := range ret.Paths {
		ret.Enable(path)
	}
	return
}

func (p *Persistor) enableAll() (ret map[string]*NginxServer) {
	ret = make(map[string]*NginxServer)
	for name := range p.servers {
		srv := p.enableServer(name)
		ret[srv.ServerName] = srv
	}
	return
}

// Disable a mapping
func (p *Persistor) Disable(name, path string) (ret map[string]*NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	if name == "" {
		return p.disableAll()
	}

	var srv *NginxServer
	if path == "" {
		srv = p.disableServer(name)
	} else {
		srv = p.disableOne(name, path)
	}

	ret = map[string]*NginxServer{
		srv.ServerName: srv,
	}
	return
}

func (p *Persistor) disableOne(name, path string) (ret *NginxServer) {
	ret = p.getServer(name)
	ret.Disable(path)
	return
}

func (p *Persistor) disableServer(name string) (ret *NginxServer) {
	ret = p.getServer(name)
	for path := range ret.Paths {
		ret.Disable(path)
	}
	return
}

func (p *Persistor) disableAll() (ret map[string]*NginxServer) {
	ret = make(map[string]*NginxServer)
	for name := range p.servers {
		srv := p.disableServer(name)
		ret[srv.ServerName] = srv
	}
	return
}

// List all server and mappings
func (p *Persistor) List() (ret map[string]*NginxServer) {
	p.Lock()
	defer p.Unlock()

	ret = make(map[string]*NginxServer)
	for _, srv := range p.servers {
		ret[srv.ServerName] = srv
	}
	return
}
