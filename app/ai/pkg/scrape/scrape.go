package scrape

import "github.com/mendableai/firecrawl-go"

var ScrapeApp *firecrawl.FirecrawlApp

func ScrapeUrl(url string) (content string, err error) {
	doc, err := ScrapeApp.ScrapeURL(url, &firecrawl.ScrapeParams{
		Formats: []string{"markdown"},
	})
	if err != nil {
		return "", err
	}

	return doc.Markdown, nil

}
