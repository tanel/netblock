package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	sectionBegin = "# BEGIN section for netblock sites"
	sectionEnd   = "# END section for netblock sites"
)

type hosts struct {
	sites []string
}

func (h *hosts) read(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	begin, end := false, false
	for _, s := range strings.Split(string(b), "\n") {
		switch s {
		case sectionBegin:
			begin = true
		case sectionEnd:
			end = true
		default:
			if begin && !end {
				cols := strings.Split(s, "\t")
				if len(cols) != 2 {
					return fmt.Errorf("2 cols expected, got %d (%s)", len(cols), s)
				}

				h.sites = append(h.sites, cols[1])
			}
		}
	}

	return nil
}

func (h *hosts) write(filename string) error {
	return nil
}

func (h *hosts) apply(cmd, site string) error {
	switch cmd {
	case cmdAdd:
		return h.add(site)
	case cmdRemove:
		return h.remove(site)
	case cmdList:
		return h.list()
	default:
		return errors.New("please specify a command: list, add, remove")
	}
}

func (h *hosts) add(site string) error {
	if site == "" {
		return errors.New("please specify a site to add")
	}

	if h.sites == nil {
		h.sites = []string{}
	}

	h.sites = append(h.sites, site)

	return nil
}

func (h *hosts) remove(site string) error {
	if site == "" {
		return errors.New("please specify a site to remove")
	}

	found := false
	for i, existingSite := range h.sites {
		if existingSite == site {
			found = true
			h.sites = append(h.sites[:i], h.sites[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("%s not found", site)
	}

	return nil
}

func (h hosts) list() error {
	for _, site := range h.sites {
		fmt.Println(site)
	}

	return nil
}
