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
	all          []string
	blockedSites []string
}

func (h *hosts) readFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	h.all = strings.Split(string(b), "\n")

	if err := h.parseBlockedHosts(); err != nil {
		return err
	}

	return nil
}

func (h *hosts) parseBlockedHosts() error {
	h.blockedSites = []string{}

	begin, end := false, false
	for _, s := range h.all {
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

				h.blockedSites = append(h.blockedSites, cols[1])
			}
		}
	}

	return nil
}

func (h *hosts) writeFile(filename string) error {
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

	if h.blockedSites == nil {
		h.blockedSites = []string{}
	}

	h.blockedSites = append(h.blockedSites, site)

	return nil
}

func (h *hosts) remove(site string) error {
	if site == "" {
		return errors.New("please specify a site to remove")
	}

	found := false
	for i, existingSite := range h.blockedSites {
		if existingSite == site {
			found = true
			h.blockedSites = append(h.blockedSites[:i], h.blockedSites[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("%s not found", site)
	}

	return nil
}

func (h hosts) list() error {
	for _, site := range h.blockedSites {
		fmt.Println(site)
	}

	return nil
}
