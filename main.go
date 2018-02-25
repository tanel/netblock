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
		if err := writeFile(filename, result); err != nil {
			return err
		}
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

type line struct {
	content       string
	isBlockedSite bool
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

type visitor func(site string, isBlocked bool)

func visit(lines []string, callback visitor) {
	begin, end := false, false
	for _, s := range lines {
		isBlocked := false
		switch s {
		case sectionBegin:
			begin = true
			end = false
		case sectionEnd:
			end = true
		default:
			isBlocked = begin && !end
		}

		callback(s, isBlocked)
	}
}

func blockedSites(lines []string) []string {
	var result []string
	visit(lines, func(s string, isBlocked bool) {
		if isBlocked {
			result = append(result, s)
		}
	})

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

func add(lines []string, site string) ([]string, error) {
	if site == "" {
		return nil, errors.New("please specify a site to add")
	}

	var result []string
	added := false
	duplicate := false
	visit(lines, func(s string, isBlocked bool) {
		newLine := "127.0.0.1\t" + site

		if isBlocked && !added {
			result = append(result, newLine)
			added = true
		}

		if s == newLine {
			duplicate = true
		}

		result = append(result, s)
	})

	if duplicate {
		return lines, nil
	}

	if !added {
		result = append(result, sectionBegin)
		result = append(result, "127.0.0.1\t"+site)
		result = append(result, sectionEnd)
	}

	return result, nil
}

func remove(lines []string, site string) ([]string, error) {
	if site == "" {
		return nil, errors.New("please specify a site to remove")
	}

	var result []string
	removed := false
	visit(lines, func(s string, isBlocked bool) {
		if isBlocked && strings.Contains(s, site) {
			removed = true
		} else {
			result = append(result, s)
		}
	})

	if !removed {
		return nil, fmt.Errorf("%s not found", site)
	}

	return result, nil
}

func list(lines []string) {
	visit(lines, func(s string, isBlocked bool) {
		if isBlocked {
			fmt.Println(host(s))
		}
	})
}
