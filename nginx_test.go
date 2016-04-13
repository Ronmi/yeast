// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import "testing"

func TestNginxExportNoname(t *testing.T) {
	expect := `server {
    client_max_body_size 250m;
    listen 80 default_server;

    location /a/ {
        proxy_pass http://a;
        include proxy_params;
        
    }

    location /b/ {
        proxy_pass http://b;
        include proxy_params;
        custom_tag 123;
    }

}`

	s := NewServer("")
	s.Create("/b/", "http://b", "custom_tag 123;")
	s.Create("/a/", "http://a", "")
	actual := s.Export()

	if actual != expect {
		t.Errorf("Noname returns %s", actual)
	}
}

func TestNginxExport(t *testing.T) {
	expect := `server {
    client_max_body_size 250m;
    server_name example.com;
    listen 80;

    location /a/ {
        proxy_pass http://a;
        include proxy_params;
        
    }

    location /b/ {
        proxy_pass http://b;
        include proxy_params;
        custom_tag 123;
    }

}`

	s := NewServer("example.com")
	s.Create("/b/", "http://b", "custom_tag 123;")
	s.Create("/a/", "http://a", "")
	actual := s.Export()

	if actual != expect {
		t.Errorf("Noname returns %s", actual)
	}
}

func TestNginxExportCustomPort(t *testing.T) {
	expect := `server {
    client_max_body_size 250m;
    server_name example.com;
    listen 81;

    location /a/ {
        proxy_pass http://a;
        include proxy_params;
        
    }

    location /b/ {
        proxy_pass http://b;
        include proxy_params;
        custom_tag 123;
    }

}`

	s := NewServer("example.com:81")
	s.Create("/b/", "http://b", "custom_tag 123;")
	s.Create("/a/", "http://a", "")
	actual := s.Export()

	if actual != expect {
		t.Errorf("Noname returns %s", actual)
	}
}
