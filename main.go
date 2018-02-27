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

	cmd, site := parseArgs(args)
	result, err := apply(lines, cmd, site)
	if err != nil {
		return err
	}

	if cmd != cmdList {
		return writeFile(filename, result)
	}

	return nil
}

const (
	cmdAdd    = "add"
	cmdRemove = "remove"
	cmdList   = "list"
)

func parseArgs(args []string) (cmd, site string) {
	switch len(args) {
	case 1:
		return args[0], ""
	case 2:
		return args[0], args[1]
	default:
		return "", ""
	}
}

const (
	sectionBegin = "# BEGIN section for netblock sites"
	sectionEnd   = "# END section for netblock sites"
)

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

func apply(lines []string, cmd, site string) ([]string, error) {
	switch cmd {
	case cmdAdd:
		return add(lines, site)
	case cmdRemove:
		return remove(lines, site)
	case cmdList:
		list(lines)
		return lines, nil
	default:
		return nil, errors.New("please specify a command: list, add, remove")
	}
}

const localhost = "127.0.0.1"

func add(lines []string, site string) ([]string, error) {
	if site == "" {
		return nil, errors.New("please specify a site to add")
	}

	// Find if already exists
	for _, s := range lines {
		if s == localhost+"\t"+site {
			return lines, nil
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

	return result, nil
}

func remove(lines []string, site string) ([]string, error) {
	if site == "" {
		return nil, errors.New("please specify a site to remove")
	}

	var result []string
	blockedSection := false
	removed := false
	for _, s := range lines {
		if s == sectionBegin {
			blockedSection = true
		} else if s == sectionEnd {
			blockedSection = false
		}

		if blockedSection && strings.Contains(s, site) {
			removed = true
			continue
		}

		result = append(result, s)
	}

	if !removed {
		return nil, fmt.Errorf("%s not found", site)
	}

	return result, nil
}

func list(lines []string) {
	for _, s := range blockedSites(lines) {
		fmt.Println(host(s))
	}
}
