package main

import (
	"flag"
	"os"
	"testing"
)

func TestParseCLIArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected struct {
			start   string
			end     string
			name    string
			keyword string
		}
	}{
		{
			args: []string{"camera", "--start", "2020-01-01", "--end", "2020-12-31", "--name", "folder-name"},
			expected: struct {
				start   string
				end     string
				name    string
				keyword string
			}{
				start:   "2020-01-01",
				end:     "2020-12-31",
				name:    "folder-name",
				keyword: "camera",
			},
		},
		{
			args: []string{"whatsapp", "--start", "2020-01-01", "--name", "folder-name"},
			expected: struct {
				start   string
				end     string
				name    string
				keyword string
			}{
				start:   "2020-01-01",
				end:     "",
				name:    "folder-name",
				keyword: "whatsapp",
			},
		},
		{
			args: []string{"camera", "--start", "2020-01-01"},
			expected: struct {
				start   string
				end     string
				name    string
				keyword string
			}{
				start:   "2020-01-01",
				end:     "",
				name:    "",
				keyword: "camera",
			},
		},
	}

	for _, test := range tests {
		// Reset the command-line flags for each test case
		flag.CommandLine = flag.NewFlagSet(test.args[0], flag.ExitOnError)
		os.Args = append([]string{"cmd"}, test.args...)

		start, end, name, keyword := parseCLIArgs()

		if start != test.expected.start {
			t.Errorf("Expected start %s, got %s", test.expected.start, start)
		}
		if end != test.expected.end {
			t.Errorf("Expected end %s, got %s", test.expected.end, end)
		}
		if name != test.expected.name {
			t.Errorf("Expected name %s, got %s", test.expected.name, name)
		}
		if keyword != test.expected.keyword {
			t.Errorf("Expected keyword %s, got %s", test.expected.keyword, keyword)
		}
	}
}
