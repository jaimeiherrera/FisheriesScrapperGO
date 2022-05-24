package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Fish struct {
	Species string `json:"species"`
	Status  string `json:"status"`
	Year    string `json:"year"`
	Region  string `json:"region"`
}

func WebScrapper() {
	spaces := regexp.MustCompile(`\s+`)
	fishes := []Fish{}

	c := colly.NewCollector(
		colly.AllowedDomains("fisheries.noaa.gov", "www.fisheries.noaa.gov"),
	)

	c.OnHTML("div.species-directory__species--8col", func(e *colly.HTMLElement) {
		species := e.DOM
		fish := Fish{
			Species: species.Find("div.species-directory__species-title--name").Text(),
			Status:  spaces.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-status").Text()), ""),
			Year:    spaces.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-year").Text()), ""),
			Region:  spaces.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-region").Text()), ""),
		}

		fishes = append(fishes, fish)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for i := 1; i < 5; i++ {
		c.Visit("https://www.fisheries.noaa.gov/species-directory/threatened-endangered?title=&species_category=any&species_status=any&regions=all&items_per_page=25&page=" + strconv.Itoa(i) + "&sort=")
	}

	writeJSON(fishes)
}

func writeJSON(fishes []Fish) {
	f, err := json.MarshalIndent(fishes, "", " ")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(f))
}

func main() {
	WebScrapper()
}
