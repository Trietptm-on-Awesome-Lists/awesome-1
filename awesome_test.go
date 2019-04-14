package awesome

import (
	"io/ioutil"
	"testing"
)

type searchTest struct {
	query string
	url   string
}

func TestSearch(t *testing.T) {
	searchTests := []searchTest{
		{"Go", "https://github.com/avelino/awesome-go#readme"},
		{"GO", "https://github.com/avelino/awesome-go#readme"},
		{"go", "https://github.com/avelino/awesome-go#readme"},
		{"Gô", ""},
	}

	b, err := ioutil.ReadFile("testdata/awesome.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range searchTests {
		r, err := Search(b, SearchRequest{Query: tc.query})
		if err != nil {
			t.Fatal(err)
		}

		for _, repo := range r.Repositories {
			if repo.Url != tc.url {
				t.Errorf("repo.Url = %q; want %q", repo.Url, tc.url)
			}
		}
	}
}
