// This file is part of Yeast
// Yeast is free software: see LICENSE.txt for more details.

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// create persistor
func cp(t *testing.T) (ret *Persistor) {
	f, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("Cannot create db when creating persistor: %s", err)
	}
	defer f.Close()
	fn := f.Name()

	f, err = ioutil.TempFile("", "nginx")
	if err != nil {
		t.Fatalf("Cannot create conf when creating persistor: %s", err)
	}
	defer f.Close()
	nginx := f.Name()

	ret = NewPersistor(fn, nginx)
	return
}

// delete persistor
func dp(p *Persistor) {
	os.Remove(p.filename)
	os.Remove(p.conffile)
}

func TestListEmpty(t *testing.T) {
	p := cp(t)
	defer dp(p)

	data := p.List()
	if len(data) != 0 {
		t.Error("Not returning empty array when no data!")
	}
}

func TestCreateAndList(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test/", "http://upstream", "")
	data := p.List()

	if len(data) != 1 {
		t.Fatalf("Created only one entry, but got %d", len(data))
	}

	actual, ok := data["test.server"]
	if !ok {
		t.Errorf("Server name differs from just created one")
	}

	mapping, ok := actual.Paths["/test/"]
	if !ok {
		t.Errorf("Path just created not found, dumping paths %#v", actual.Paths)
	}

	if mapping.Upstream != "http://upstream" {
		t.Errorf("Upstream %s differs from just created one", mapping.Upstream)
	}

	if mapping.CustomTags != "" {
		t.Errorf("Custom tag %s differs from just created one", mapping.CustomTags)
	}

	if !mapping.Enabled {
		t.Error("Newly created data should be auto-enabled")
	}
}

func TestModify(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test/", "http://upstream", "")
	p.Modify("test.server", "/test/", "/orz/", "http://orz", "")
	data := p.List()
	if _, ok := data["test.server"].Paths["/orz/"]; !ok {
		t.Error("Cannot find modified path")
	}
}

func TestDelete(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test/", "http://upstream", "")
	p.Delete("test.server", "/test/")

	data := p.List()
	if len(data) != 0 {
		t.Error("Not returning empty array after deleting just created data")
	}
}

func TestDisableOne(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test/", "http://upstream", "")
	p.Disable("test.server", "/test/")

	data := p.List()
	if data["test.server"].Paths["/test/"].Enabled {
		t.Error("Disable is not working")
	}
}

func TestEnableOne(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test/", "http://upstream", "")
	p.Disable("test.server", "/test/")
	p.Enable("test.server", "/test/")

	data := p.List()
	if !data["test.server"].Paths["/test/"].Enabled {
		t.Error("Enable is not working")
	}
}

func TestDisableServer(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test1/", "http://upstream", "")
	p.Create("test.server", "/test2/", "http://upstream", "")
	p.Disable("test.server", "")

	data := p.List()
	for _, path := range []string{"/test1/", "/test2/"} {
		if data["test.server"].Paths[path].Enabled {
			t.Error("Disable server is not working")
		}
	}
}

func TestEnableServer(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test.server", "/test1/", "http://upstream", "")
	p.Create("test.server", "/test2/", "http://upstream", "")
	p.Disable("test.server", "")
	p.Enable("test.server", "")

	data := p.List()
	for _, path := range []string{"/test1/", "/test2/"} {
		if !data["test.server"].Paths[path].Enabled {
			t.Error("Enable server is not working")
		}
	}
}

func TestDisableAll(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test1.server", "/test1/", "http://upstream", "")
	p.Create("test1.server", "/test2/", "http://upstream", "")
	p.Create("test2.server", "/test1/", "http://upstream", "")
	p.Create("test2.server", "/test2/", "http://upstream", "")
	p.Disable("", "")

	data := p.List()
	for _, data := range data {
		for _, path := range []string{"/test1/", "/test2/"} {
			if data.Paths[path].Enabled {
				t.Error("Disable all server is not working")
			}
		}
	}
}

func TestEnableAll(t *testing.T) {
	p := cp(t)
	defer dp(p)

	p.Create("test1.server", "/test1/", "http://upstream", "")
	p.Create("test1.server", "/test2/", "http://upstream", "")
	p.Create("test2.server", "/test1/", "http://upstream", "")
	p.Create("test2.server", "/test2/", "http://upstream", "")
	p.Disable("", "")
	p.Enable("", "")

	data := p.List()
	for _, data := range data {
		for _, path := range []string{"/test1/", "/test2/"} {
			if !data.Paths[path].Enabled {
				t.Error("Enable all server is not working")
			}
		}
	}
}
