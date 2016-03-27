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

// Set path to upstream mapping
func (p *Persistor) Set(name, path, upstream, custom string) (ret *NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	return p.getServer(name).Set(path, upstream, custom)
}

// Unset a path-upstream mapping
func (p *Persistor) Unset(name, path string) (ok bool) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	s := p.getServer(name)
	s.Unset(path)

	if s.Len() < 1 {
		delete(p.servers, name)
		ok = true
	}
	return
}

// Enable a mapping
func (p *Persistor) Enable(name, path string) (ret []*NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	if name == "" {
		return p.enableAll()
	}

	if path == "" {
		return []*NginxServer{p.enableServer(name)}
	}

	return []*NginxServer{p.enableOne(name, path)}
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

func (p *Persistor) enableAll() (ret []*NginxServer) {
	ret = make([]*NginxServer, 0, len(p.servers))
	for name := range p.servers {
		ret = append(ret, p.enableServer(name))
	}
	return
}

// Disable a mapping
func (p *Persistor) Disable(name, path string) (ret []*NginxServer) {
	p.Lock()
	defer p.Unlock()
	defer p.doSave()

	if name == "" {
		return p.disableAll()
	}

	if path == "" {
		return []*NginxServer{p.disableServer(name)}
	}

	return []*NginxServer{p.disableOne(name, path)}
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

func (p *Persistor) disableAll() (ret []*NginxServer) {
	ret = make([]*NginxServer, 0, len(p.servers))
	for name := range p.servers {
		ret = append(ret, p.disableServer(name))
	}
	return
}

// List all server and mappings
func (p *Persistor) List() (ret map[string]map[string]*Mapping) {
	ret = map[string]map[string]*Mapping{}
	p.Lock()
	defer p.Unlock()

	for name, srv := range p.servers {
		ret[name] = srv.List()
	}

	return
}
