package awesome

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday/v2"
)

func Search(b []byte, r SearchRequest) (SearchResponse, error) {
	output := blackfriday.Run(b)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(output)))
	if err != nil {
		return SearchResponse{}, err
	}

	var repos []*Repository
	doc.Find("ul li a[href*='https://github.com/']").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.ToLower(r.Query) != strings.ToLower(text) {
			return
		}

		url, ok := s.Attr("href")
		if !ok {
			return
		}
		repos = append(repos, &Repository{Url: url})
	})

	return SearchResponse{Repositories: repos}, nil
}
