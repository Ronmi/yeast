package main

import (
	"sort"
	"strings"
	"sync"
)

type nginxConf struct {
	data []string
}

func newConf() *nginxConf {
	return &nginxConf{make([]string, 0, 2)}
}

func (c *nginxConf) add(line string) *nginxConf {
	c.data = append(c.data, line)
	return c
}

func (c *nginxConf) indent(line string, level int) *nginxConf {
	return c.add(strings.Repeat("    ", level) + line)
}

func (c *nginxConf) export() string {
	return strings.Join(c.data, "\n")
}

// Mapping is base structure of path-upstream mapping
type Mapping struct {
	Upstream string `json:"upstream"`
	Enabled  bool   `json:"-"`
}

// NginxServer represents server segment of nginx conf
type NginxServer struct {
	ServerName   string              `json:"name"`
	Paths        map[string]*Mapping `json:"paths"`
	length       int
	sync.RWMutex `json:"-"`
}

// NewServer creates a new NginxServer
func NewServer(name string) *NginxServer {
	return &NginxServer{
		name,
		map[string]*Mapping{},
		0,
		sync.RWMutex{},
	}
}

// Set path => upstream mapping
func (s *NginxServer) Set(path, upstream string) {
	s.Lock()
	defer s.Unlock()

	s.Paths[path] = &Mapping{upstream, true}
	s.length++
}

// Unset path => upstream mapping
func (s *NginxServer) Unset(path string) {
	s.Lock()
	defer s.Unlock()

	delete(s.Paths, path)
	s.length--
}

// Disable a mapping, but not deleting it. noop if not found
func (s *NginxServer) Disable(path string) {
	s.Lock()
	defer s.Unlock()

	if mapping, ok := s.Paths[path]; ok {
		mapping.Enabled = false
	}
}

// Enable a mapping. noop if not found
func (s *NginxServer) Enable(path string) {
	s.Lock()
	defer s.Unlock()

	if mapping, ok := s.Paths[path]; ok {
		mapping.Enabled = true
	}
}

// List all mapping data
func (s *NginxServer) List() (ret map[string]*Mapping) {
	ret = map[string]*Mapping{}
	s.RLock()
	defer s.RUnlock()

	for path, mapping := range s.Paths {
		ret[path] = mapping
	}

	return
}

// Len returns how many mapping data in this server
func (s *NginxServer) Len() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.Paths)
}

// Export to string
func (s *NginxServer) Export() string {
	s.RLock()
	defer s.RUnlock()

	ret := newConf()
	ret.add("server {")
	ret.indent("client_max_body_size 250m;", 1)

	if s.ServerName == "" {
		ret.indent("listen 80 default_server;", 1)
	} else {
		ret.indent("server_name "+s.ServerName+";", 1)
	}
	ret.add("")

	buf := make([]string, 0, len(s.Paths))

	for path := range s.Paths {
		buf = append(buf, path)
	}
	sort.Strings(buf)

	for _, path := range buf {
		mapping := s.Paths[path]
		if !mapping.Enabled {
			continue
		}
		ret.indent("location "+path+" {", 1)
		ret.indent("proxy_pass "+mapping.Upstream+";", 2)
		ret.indent("include proxy_params;", 2)
		ret.indent("}", 1)
		ret.add("")
	}

	ret.add("}")
	return ret.export()
}
