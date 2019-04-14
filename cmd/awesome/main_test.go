package main

import (
	"bytes"
	"strings"
	"testing"
)

type searchTest struct {
	in  string
	out string
}

func TestSearch(t *testing.T) {
	searchTests := []searchTest{
		{"Go", "https://github.com/avelino/awesome-go#readme"},
		{"GO", "https://github.com/avelino/awesome-go#readme"},
		{"go", "https://github.com/avelino/awesome-go#readme"},
		{"Gô", ""},
	}

	for _, tc := range searchTests {
		out := &bytes.Buffer{}
		c := &cli{&bytes.Buffer{}, out, &bytes.Buffer{}}
		args := []string{"awesome", tc.in}
		exitCode := c.run(args)
		if exitCode != exitCodeSuccess {
			t.Errorf("%q exit status %d; want %d", args, exitCode, exitCodeSuccess)
		}
		if strings.TrimSpace(out.String()) != tc.out {
			t.Errorf("%q got %q; want %q", args, strings.TrimSpace(strings.TrimSpace(out.String())), tc.out)
		}
	}
}
