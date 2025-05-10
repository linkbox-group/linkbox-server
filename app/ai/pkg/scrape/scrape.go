package scrape

import (
	"github.com/mendableai/firecrawl-go"
	"github.com/spf13/viper"
	"log"
)

var App *firecrawl.FirecrawlApp

type FireCrawApp struct {
	*firecrawl.FirecrawlApp
}

func NewApp() *FireCrawApp {

	app, err := firecrawl.NewFirecrawlApp(viper.GetString("firecrawl.apikey"), viper.GetString("firecrawl.apiurl"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &FireCrawApp{app}
}

func (a *FireCrawApp) ExtractContent(url string) (content string, err error) {
	doc, err := a.ScrapeURL(url, &firecrawl.ScrapeParams{
		Formats: []string{"markdown"},
	})
	if err != nil {
		return "", err
	}

	return doc.Markdown, nil

}
