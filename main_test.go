package main

import (
	"testing"
)

func Test_parse(t *testing.T) {
	var examples = []struct {
		args         []string
		expectedCmd  string
		expectedSite string
	}{
		{
			args:         nil,
			expectedCmd:  "",
			expectedSite: "",
		},
		{
			args:         []string{cmdList},
			expectedCmd:  cmdList,
			expectedSite: "",
		},
		{
			args:         []string{cmdAdd, "test.com"},
			expectedCmd:  cmdAdd,
			expectedSite: "test.com",
		},
		{
			args:         []string{cmdRemove, "test.com"},
			expectedCmd:  cmdRemove,
			expectedSite: "test.com",
		},
	}

	for _, example := range examples {
		cmd, site := parse(example.args)
		if cmd != example.expectedCmd {
			t.Errorf("cmd expected %s but got %s", example.expectedCmd, cmd)
		}

		if site != example.expectedSite {
			t.Errorf("site expected %s but got %s", example.expectedSite, site)
		}
	}
}
