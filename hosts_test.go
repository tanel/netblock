package main

import (
	"path/filepath"
	"testing"
)

func Test_hosts_apply(t *testing.T) {
	var h hosts
	if err := h.apply("", ""); err == nil {
		t.Error("error expected")
	}
}

func Test_hosts_apply_add_site(t *testing.T) {
	var h hosts
	if err := h.apply(cmdAdd, "test.com"); err != nil {
		t.Error(err)
	}

	found := false
	for _, site := range h.blockedSites {
		if site == "test.com" {
			found = true
			break
		}
	}

	if !found {
		t.Error("site not added")
	}
}

func Test_hosts_apply_remove_site_nonexistant(t *testing.T) {
	var h hosts
	if err := h.apply(cmdRemove, "test.com"); err == nil {
		t.Error("error expected")
	}
}

func Test_hosts_apply_remove_site(t *testing.T) {
	var h hosts
	h.add("test.com")
	if err := h.apply(cmdRemove, "test.com"); err != nil {
		t.Error(err)
	}

	found := false
	for _, site := range h.blockedSites {
		if site == "test.com" {
			found = true
			break
		}
	}

	if found {
		t.Error("site not remove")
	}

}

func Test_hosts_apply_list(t *testing.T) {
	var h hosts
	if err := h.list(); err != nil {
		t.Error(err)
	}
}

func Test_readFile_nosites(t *testing.T) {
	var h hosts
	if err := h.readFile(filepath.Join("testdata", "no-sites")); err != nil {
		t.Error(err)
	}

	if len(h.blockedSites) != 0 {
		t.Errorf("0 sites expected, got %d", len(h.blockedSites))
	}
}

func Test_readFile(t *testing.T) {
	var h hosts
	if err := h.readFile(filepath.Join("testdata", "2-sites")); err != nil {
		t.Error(err)
	}

	if len(h.blockedSites) != 2 {
		t.Errorf("2 sites expected, got %d", len(h.blockedSites))
	}

	if h.blockedSites[0] != "test.com" {
		t.Errorf("unexpected site %s", h.blockedSites[0])
	}

	if h.blockedSites[1] != "www.test.com" {
		t.Errorf("unexpected site %s", h.blockedSites[1])
	}
}
