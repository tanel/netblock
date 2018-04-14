package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if err := run("/etc/hosts", os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(filename string, args []string) error {
	lines, err := readFile(filename)
	if err != nil {
		return err
	}

	cmd, sites := parseArgs(args)
	result, err := apply(lines, cmd, sites)
	if err != nil {
		return err
	}

	if cmd != cmdList {
		return writeFile(filename, result)
	}

	return nil
}

func parseArgs(args []string) (cmd string, sites []string) {
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		return args[0], nil
	default:
		return args[0], args[1:]
	}
}

func host(s string) string {
	cols := strings.Split(s, "\t")
	if len(cols) == 2 {
		return cols[1]
	}

	return ""
}

func readFile(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

func writeFile(filename string, lines []string) error {
	s := strings.Join(lines, "\n")
	return ioutil.WriteFile(filename, []byte(s), 0644)
}

const (
	sectionBegin = "# BEGIN section for netblock sites"
	sectionEnd   = "# END section for netblock sites"
)

func blockedSites(lines []string) []string {
	var result []string
	blockedSection := false
	for _, s := range lines {
		if s == sectionBegin {
			blockedSection = true
		} else if s == sectionEnd {
			blockedSection = false
		} else if blockedSection {
			result = append(result, s)
		}
	}

	return result
}

const (
	cmdAdd    = "add"
	cmdRemove = "remove"
	cmdList   = "list"
)

func apply(lines []string, cmd string, sites []string) ([]string, error) {
	switch cmd {
	case cmdAdd:
		return add(lines, sites)
	case cmdRemove:
		return remove(lines, sites)
	case cmdList:
		list(lines)
		return lines, nil
	default:
		return nil, errors.New("please specify a command: list, add, remove")
	}
}

func add(lines []string, sites []string) ([]string, error) {
	if len(sites) == 0 {
		return nil, errors.New("please specify site(s) to add")
	}

	for _, site := range sites {
		lines = addSite(lines, site)
		if !strings.HasPrefix(site, "www.") {
			lines = addSite(lines, "www."+site)
		}
	}

	return lines, nil
}

const localhost = "0.0.0.0"

func addSite(lines []string, site string) []string {
	// Find if already exists
	for _, s := range lines {
		if s == localhost+"\t"+site {
			return lines
		}
	}

	// Find a place to insert
	var result []string
	added := false
	for _, s := range lines {
		result = append(result, s)

		if s == sectionBegin {
			result = append(result, localhost+"\t"+site)
			added = true
		}
	}

	// Append if not inserted above
	if !added {
		result = append(result, sectionBegin)
		result = append(result, localhost+"\t"+site)
		result = append(result, sectionEnd)
	}

	return result
}

func remove(lines []string, sites []string) ([]string, error) {
	if len(sites) == 0 {
		return nil, errors.New("please specify site(s) to remove")
	}

	for _, site := range sites {
		var removed bool
		lines, removed = removeSite(lines, site)
		if !removed {
			return nil, fmt.Errorf("%s not found", site)
		}

		if !strings.HasPrefix(site, "www.") {
			lines, _ = removeSite(lines, "www."+site)
		}
	}

	return lines, nil
}

func removeSite(lines []string, site string) (result []string, removed bool) {
	blockedSection := false
	for _, s := range lines {
		if s == sectionBegin {
			blockedSection = true
		} else if s == sectionEnd {
			blockedSection = false
		}

		if blockedSection && s == localhost+"\t"+site {
			removed = true
			continue
		}

		result = append(result, s)
	}

	return result, removed
}

func list(lines []string) {
	for _, s := range blockedSites(lines) {
		fmt.Println(host(s))
	}
}
