package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gocolly/colly"
)

type RSS struct {
	Version     string `json:"version"`
	UserComment string `json:"user_comment"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HomePageUrl string `json:"home_page_url"`
	FeedUrl     string `json:"feed_url"`
	Author      Author `json:"author"`
	Items       []Episode `json:"items"`
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Episode struct {
	Title     string
	Lead      string
	Url       string
	Image     string
	Published string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	collector := colly.NewCollector(
		colly.AllowedDomains("sic.pt"),
		colly.CacheDir("./cache"),
	)

	episodes := make([]Episode, 0, 200)

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Requesting:", r.URL.String())
	})

	collector.OnHTML("article", func(e *colly.HTMLElement) {
		episode := Episode{
			Title: e.ChildText(".textDetails .title a"),
			Lead: e.ChildText(".textDetails .lead"),
			Url: "https://sic.pt/" + e.ChildAttr("a", "href"),
			Image: "https:" + e.ChildAttr("img", "src"),
			Published: e.ChildText(".publishedDate"),
		}

		episodes = append(episodes, episode)
	})

	collector.Visit("https://sic.pt/Programas/governo-sombra/videos")

	generateFile(episodes)
}

func generateFile(episodes []Episode) {
	domain := os.Getenv("DOMAIN")
	rss := RSS{
		domain + "/go-verno-sombra/version/1",
		"Scraper in Go to generate a json feed for Governo Sombra",
		"Governo Sombra",
		"SIC",
		domain,
		domain + "/go-verno-sombra/feed.json",
		Author{
			"Nuno Lopes",
			domain,
		},
		episodes,
	}

	json, _ := json.MarshalIndent(rss, "", " ")
	err := ioutil.WriteFile("./feed.json", json, 0644)

	if err != nil {
		log.Println(err)

		return
	}

	log.Println("File './feed.json' was successfully generated")
}