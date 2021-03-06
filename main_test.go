package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	var examples = []struct {
		args          []string
		expectedCmd   string
		expectedSites []string
	}{
		{
			args:          nil,
			expectedCmd:   "",
			expectedSites: nil,
		},
		{
			args:          []string{cmdList},
			expectedCmd:   cmdList,
			expectedSites: nil,
		},
		{
			args:          []string{cmdAdd, "test.com"},
			expectedCmd:   cmdAdd,
			expectedSites: []string{"test.com"},
		},
		{
			args:          []string{cmdAdd, "test.com", "something.eu"},
			expectedCmd:   cmdAdd,
			expectedSites: []string{"test.com", "something.eu"},
		},
		{
			args:          []string{cmdRemove, "test.com"},
			expectedCmd:   cmdRemove,
			expectedSites: []string{"test.com"},
		},
		{
			args:          []string{cmdRemove, "test.com", "something.eu"},
			expectedCmd:   cmdRemove,
			expectedSites: []string{"test.com", "something.eu"},
		},
		{
			args:          []string{cmdEnable},
			expectedCmd:   cmdEnable,
			expectedSites: nil,
		},
		{
			args:          []string{cmdDisable},
			expectedCmd:   cmdDisable,
			expectedSites: nil,
		},
	}

	for _, example := range examples {
		cmd, sites := parseArgs(example.args)
		if cmd != example.expectedCmd {
			t.Errorf("cmd expected %s but got %s", example.expectedCmd, cmd)
		}

		if len(sites) != len(example.expectedSites) {
			t.Errorf("%d site(s) expected but got %d", len(example.expectedSites), len(sites))
		}
	}
}

func Test_readFile_NoBlockedSites(t *testing.T) {
	result, err := readFile(filepath.Join("testdata", "no-sites"))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 10 {
		t.Errorf("10 lines expected, got %d", len(result))
	}

	blocked := blockedSites(result)
	if len(blocked) != 0 {
		t.Errorf("no sites expected but got %d", len(blocked))
	}
}

func Test_readFile_NonExistantFile(t *testing.T) {
	if _, err := readFile("does-not-exist"); err == nil {
		t.Error("error expected")
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

	sites := blockedSites(result)

	if len(sites) != 2 {
		t.Errorf("2 sites expected, got %d", len(sites))
	}

	if sites[0] != localhost+"\ttest.com" {
		t.Errorf("unexpected site %s", sites[0])
	}

	if sites[1] != localhost+"\twww.test.com" {
		t.Errorf("unexpected site %s", sites[1])
	}
}

func Test_host(t *testing.T) {
	result := host(localhost + "\twww.test.com")
	if result != "www.test.com" {
		t.Errorf("expected www.test.com, got %s", result)
	}
}

func Test_host_InvalidLine(t *testing.T) {
	result := host(localhost + " www.test.com")
	if result != "" {
		t.Errorf("expected no result, got %s", result)
	}
}

func copy(source, destination string) error {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, b, 0644)
}

func Test_run_InvalidFile(t *testing.T) {
	if err := run("this file does not exist", nil); err == nil {
		t.Error("error expected")
	}
}

func Test_run_RemoveWithoutArg(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove}); err == nil {
		t.Error("error expected")
	}
}

func Test_run_Remove(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove, "test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 0 {
		t.Errorf("0 sites expected after removal, got %d", len(blocked))
	}
}

func Test_run_RemoveMultipleSites(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "another.com"}); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove, "test.com", "another.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 0 {
		t.Errorf("0 sites expected after removal, got %d", len(blocked))
	}
}

func Test_run_RemoveRemovesWWW(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 sites expected, got %d", len(blocked))
	}

	if err := run(input, []string{cmdRemove, "test.com"}); err != nil {
		t.Error(err)
	}

	result, err = readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked = blockedSites(result)
	if len(blocked) != 0 {
		t.Errorf("0 sites expected after removal, got %d", len(blocked))
	}
}

func Test_run_RemoveNonExistant(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove, "test.com"}); err == nil {
		t.Error("error expected")
	}
}

func Test_run_MoreSpecificRemoval(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove, "www.test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 1 {
		t.Errorf("1 site expected after removal, got %d", len(blocked))
	}

	if blocked[0] != localhost+"\ttest.com" {
		t.Errorf("unexpected result %s", blocked[0])
	}
}

func Test_run_AddToEmptyFile(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 1 {
		t.Errorf("1 site expected after adding, got %d", len(blocked))
	}

	if blocked[0] != localhost+"\twww.test.com" {
		t.Errorf("unexpected result %s", blocked[0])
	}
}

func Test_run_AddMultipleSites(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com", "www.something.eu"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 sites expected after adding, got %d", len(blocked))
	}

	if blocked[0] != localhost+"\twww.something.eu" {
		t.Errorf("unexpected result %s", blocked[0])
	}

	if blocked[1] != localhost+"\twww.test.com" {
		t.Errorf("unexpected result %s", blocked[1])
	}
}

func Test_run_AddAddsWWW(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 sites expected after adding, got %d", len(blocked))
	}

	if blocked[0] != localhost+"\twww.test.com" {
		t.Errorf("unexpected result %s", blocked[0])
	}

	if blocked[1] != localhost+"\ttest.com" {
		t.Errorf("unexpected result %s", blocked[1])
	}
}

func Test_run_AddMultipleTimes(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	sections := 0
	for _, s := range result {
		if s == sectionBegin {
			sections++
		}
	}

	if sections != 1 {
		t.Errorf("1 section expected, got %d", sections)
	}

	blocked := blockedSites(result)
	if len(blocked) != 1 {
		t.Errorf("1 site expected after adding, got %d", len(blocked))
	}
}

func Test_run_AddRemoveAddCreatesOneSectionOnly(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdRemove, "www.test.com"}); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	sections := 0
	for _, s := range result {
		if s == sectionBegin {
			sections++
		}
	}

	if sections != 1 {
		t.Errorf("1 section expected, got %d", sections)
	}

	blocked := blockedSites(result)
	if len(blocked) != 1 {
		t.Errorf("1 site expected after adding, got %d", len(blocked))
	}
}

func Test_run_AddToNonEmptyFile(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd, "www.test.com"}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 site expected after adding existing site, got %d", len(blocked))
	}
}

func Test_run_List(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdList}); err != nil {
		t.Error(err)
	}
}

func Test_run_ListNoBlockedSites(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdList}); err != nil {
		t.Error(err)
	}
}

func Test_run_WithoutArgs(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{}); err == nil {
		t.Error("error expected")
	}
}

func Test_run_AddWithoutArg(t *testing.T) {
	example := filepath.Join("testdata", "no-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdAdd}); err == nil {
		t.Error("error expected")
	}
}

func Test_run_Disable(t *testing.T) {
	example := filepath.Join("testdata", "2-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdDisable}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 blocked sites expected, got %d", len(blocked))
		t.FailNow()
	}

	expected1 := "#0.0.0.0	test.com"
	if blocked[0] != expected1 {
		t.Errorf("expected '%s', got '%s'", expected1, blocked[0])
	}

	expected2 := "#0.0.0.0	www.test.com"
	if blocked[1] != expected2 {
		t.Errorf("expected '%s', got '%s'", expected2, blocked[1])
	}

}

func Test_run_Enable(t *testing.T) {
	example := filepath.Join("testdata", "2-disabled-sites")
	input := filepath.Join("testdata", "output", "testfile")
	if err := copy(example, input); err != nil {
		t.Error(err)
	}

	if err := run(input, []string{cmdEnable}); err != nil {
		t.Error(err)
	}

	result, err := readFile(input)
	if err != nil {
		t.Error(err)
	}

	blocked := blockedSites(result)
	if len(blocked) != 2 {
		t.Errorf("2 blocked sites expected, got %d", len(blocked))
		t.FailNow()
	}

	expected1 := "0.0.0.0	test.com"
	if blocked[0] != expected1 {
		t.Errorf("expected '%s', got '%s'", expected1, blocked[0])
	}

	expected2 := "0.0.0.0	www.test.com"
	if blocked[1] != expected2 {
		t.Errorf("expected '%s', got '%s'", expected2, blocked[1])
	}

}
