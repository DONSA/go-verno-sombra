package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type RSS struct {
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HomePageUrl string `json:"home_page_url"`
	FeedUrl     string `json:"feed_url"`
	Author      Author `json:"author"`
	Icon        string `json:"icon"`
	Favicon     string `json:"favicon"`
	Items       []Episode `json:"items"`
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Episode struct {
	Title         string `json:"title"`
	DatePublished string `json:"date_published"`
	Id            string `json:"id"`
	Url           string `json:"url"`
	ContentHtml   string `json:"content_html"`
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
			//DatePublished: e.ChildText(".publishedDate"),
			DatePublished: time.Now().Format(time.RFC3339),
			Id: "https://sic.pt/" + e.ChildAttr("a", "href"),
			Url: "https://sic.pt/" + e.ChildAttr("a", "href"),
			ContentHtml: e.ChildText(".textDetails .lead"),

		}

		episodes = append(episodes, episode)
	})

	collector.Visit("https://sic.pt/Programas/governo-sombra/videos")

	generateFile(episodes)
}

func generateFile(episodes []Episode) {
	domain := os.Getenv("DOMAIN")
	rss := RSS{
		"https://jsonfeed.org/version/1",
		"Governo Sombra",
		"Scraper in Go to generate a json feed",
		"https://sic.pt/Programas/governo-sombra/videos",
		domain + "/go-verno-sombra/feeds/json",
		Author{
			"Nuno Lopes",
			domain,
		},
		"https://static.impresa.pt/sic/2039//assets/gfx/icon.png",
		"https://sic.pt/favicon.ico?v=2",
		episodes,
	}

	json, _ := json.MarshalIndent(rss, "", " ")
	err := ioutil.WriteFile("./feeds/json", json, 0644)

	if err != nil {
		log.Println(err)

		return
	}

	log.Println("File './feeds/json' was successfully generated")
}