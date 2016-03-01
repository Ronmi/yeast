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
	s.Set("/b/", "http://b", "custom_tag 123;")
	s.Set("/a/", "http://a", "")
	actual := s.Export()

	if actual != expect {
		t.Errorf("Noname returns %s", actual)
	}
}

func TestNginxExport(t *testing.T) {
	expect := `server {
    client_max_body_size 250m;
    server_name example.com;

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
	s.Set("/b/", "http://b", "custom_tag 123;")
	s.Set("/a/", "http://a", "")
	actual := s.Export()

	if actual != expect {
		t.Errorf("Noname returns %s", actual)
	}
}
