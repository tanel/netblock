package main

import (
	"path/filepath"
	"testing"
)

func Test_parseArgs(t *testing.T) {
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
		cmd, site := parseArgs(example.args)
		if cmd != example.expectedCmd {
			t.Errorf("cmd expected %s but got %s", example.expectedCmd, cmd)
		}

		if site != example.expectedSite {
			t.Errorf("site expected %s but got %s", example.expectedSite, site)
		}
	}
}

func Test_apply(t *testing.T) {
	var lines []line
	result, err := apply(lines, "", "")
	if err == nil {
		t.Error("error expected")
	}

	if result != nil {
		t.Error("nil result expected")
	}
}

func Test_apply_add_site(t *testing.T) {
	var lines []line
	result, err := apply(lines, cmdAdd, "test.com")
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, site := range result {
		if site.content == "127.0.0.1\ttest.com" {
			found = true
			break
		}
	}

	if !found {
		t.Error("site not added")
	}
}

func Test_apply_remove_site_nonexistant(t *testing.T) {
	var lines []line
	result, err := apply(lines, cmdRemove, "test.com")
	if err == nil {
		t.Error("error expected")
	}

	if result != nil {
		t.Error("nil result expected")
	}
}

func Test_apply_remove_site(t *testing.T) {
	var lines []line
	lines = append(lines, line{
		content:       "test.com",
		isBlockedSite: true,
	})
	result, err := apply(lines, cmdRemove, "test.com")
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, site := range result {
		if site.content == "test.com" {
			found = true
			break
		}
	}

	if found {
		t.Error("site not remove")
	}

}

func Test_apply_list(t *testing.T) {
	var lines []line
	lines = append(lines, line{
		content: "test.com",
	})
	result, err := apply(lines, cmdList, "")
	if err != nil {
		t.Error(err)
	}

	if len(result) != len(lines) {
		t.Error("output should be same as input")
	}

	if result[0].content != lines[0].content {
		t.Error("output should be same as input")
	}

	if result[0].isBlockedSite != lines[0].isBlockedSite {
		t.Error("output should be same as input")
	}
}

func Test_readFile_nosites(t *testing.T) {
	result, err := readFile(filepath.Join("testdata", "no-sites"))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 10 {
		t.Errorf("10 lines expected, got %d", len(result))
	}

	for _, l := range result {
		if l.isBlockedSite {
			t.Error("no sites expected")
		}
	}
}

func Test_readFile(t *testing.T) {
	result, err := readFile(filepath.Join("testdata", "2-sites"))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 15 {
		t.Errorf("2 lines expected, got %d", len(result))
	}

	var sites []string
	for _, l := range result {
		if l.isBlockedSite {
			sites = append(sites, l.content)
		}
	}

	if len(sites) != 2 {
		t.Errorf("2 sites expected, got %d", len(sites))
	}

	if sites[0] != "127.0.0.1\ttest.com" {
		t.Errorf("unexpected site %s", sites[0])
	}

	if sites[1] != "127.0.0.1\twww.test.com" {
		t.Errorf("unexpected site %s", sites[1])
	}
}
