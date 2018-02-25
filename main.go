package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const filename = "/etc/hosts"

func run() error {
	lines, err := readFile(filename)
	if err != nil {
		return err
	}

	cmd, site := parseArgs(os.Args[1:])
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

func host(l line) (string, error) {
	cols := strings.Split(l.content, "\t")
	if len(cols) != 2 {
		return "", fmt.Errorf("2 cols expected, got %d (%s)", len(cols), l.content)
	}

	return cols[1], nil
}

func readFile(filename string) ([]line, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parseLines(strings.Split(string(b), "\n")), nil
}

func writeFile(filename string, lines []line) error {
	// FIXME:
	return nil
}

func parseLines(all []string) []line {
	var result []line
	begin, end := false, false
	for _, s := range all {
		var l line
		l.content = s

		switch s {
		case sectionBegin:
			begin = true
			end = false
		case sectionEnd:
			end = true
		default:
			l.isBlockedSite = begin && !end
		}

		result = append(result, l)
	}

	return result
}

func apply(lines []line, cmd, site string) ([]line, error) {
	switch cmd {
	case cmdAdd:
		return add(lines, site)
	case cmdRemove:
		return remove(lines, site)
	case cmdList:
		return list(lines)
	default:
		return nil, errors.New("please specify a command: list, add, remove")
	}
}

func add(lines []line, site string) ([]line, error) {
	if site == "" {
		return nil, errors.New("please specify a site to add")
	}

	lineToAdd := line{
		content:       "127.0.0.1\t" + site,
		isBlockedSite: true,
	}

	var result []line
	added := false
	for _, l := range lines {
		if l.isBlockedSite && !added {
			result = append(result, lineToAdd)
			added = true
		}

		result = append(result, l)
	}

	if !added {
		result = append(result, line{
			content:       sectionBegin,
			isBlockedSite: false,
		})
		result = append(result, lineToAdd)
		result = append(result, line{
			content:       sectionEnd,
			isBlockedSite: false,
		})
	}

	return result, nil
}

func remove(lines []line, site string) ([]line, error) {
	if site == "" {
		return nil, errors.New("please specify a site to remove")
	}

	var result []line
	removed := false
	for _, l := range lines {
		if l.isBlockedSite && strings.Contains(l.content, site) {
			removed = true
		} else {
			result = append(result, l)
		}
	}

	if !removed {
		return nil, fmt.Errorf("%s not found", site)
	}

	return result, nil
}

func list(lines []line) ([]line, error) {
	for _, l := range lines {
		if l.isBlockedSite {
			fmt.Println(host(l))
		}
	}

	return lines, nil
}
