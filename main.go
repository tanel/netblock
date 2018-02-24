package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var h hosts
	if err := h.read("//etc/hosts"); err != nil {
		return err
	}

	cmd, site := parse(os.Args[1:])
	if err := h.apply(cmd, site); err != nil {
		return err
	}

	return nil
}

const (
	cmdAdd    = "add"
	cmdRemove = "remove"
	cmdList   = "list"
)

func parse(args []string) (cmd, site string) {
	switch len(args) {
	case 1:
		return args[0], ""
	case 2:
		return args[0], args[1]
	default:
		return "", ""
	}
}
