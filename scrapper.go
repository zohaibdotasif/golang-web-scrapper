package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const (
	FirstPageToScrape = "https://scrapeme.live/shop"
	MaxPagesToVisit   = 5
)

// defining a data structure to store the scraped data
type PokemonProduct struct {
	url, image, name, price string
}

func main() {
	fmt.Println("GOLANG - WEB SCRAPPER")

	// initializing the slice of structs that will contain the scraped data
	var pokemonProducts []*PokemonProduct

	// to scrape a single webpage
	scrapeFirstWebPage(pokemonProducts)

	pokemonProducts = nil
	fmt.Println("Emptied the slice: ", pokemonProducts)

	// to scrape a multiple webpages
	scrapeMultiplePages(pokemonProducts)
}

func scrapeMultiplePages(pokemonProducts []*PokemonProduct) {
	pagesToScrape := []string{}

	currentPage := 1

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		newPageLink := e.Attr("href")

		if !contains(pagesToScrape, newPageLink) {
			pagesToScrape = append(pagesToScrape, newPageLink)
		}
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.price = e.ChildText(".price")
		pokemonProduct.url = e.ChildAttr("a", "href")

		pokemonProducts = append(pokemonProducts, &pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		if len(pagesToScrape) != 0 && currentPage < MaxPagesToVisit {
			pageToScrape := pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			currentPage++

			if err := c.Visit(pageToScrape); err != nil {
				panic(err)
			}
		}
	})

	if err := c.Visit(FirstPageToScrape); err != nil {
		panic(err)
	}

	exportToCSV("scraped-multiple-pages", pokemonProducts)
}

func scrapeFirstWebPage(pokemonProducts []*PokemonProduct) {
	// inititalize a colly collector
	c := colly.NewCollector()

	// attach different types of callback functions to a collector
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// initializing a new PokemonProduct instance
		pokemonProduct := PokemonProduct{}

		// scraping the data of interest
		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		// adding the product instance with scraped data to the slice of products
		pokemonProducts = append(pokemonProducts, &pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scrapped!")
	})

	// visit a webpage
	if err := c.Visit(FirstPageToScrape); err != nil {
		panic(err)
	}

	exportToCSV("scraped-single-page", pokemonProducts)
}

func exportToCSV(filename string, pokemonProducts []*PokemonProduct) error {
	// opening the CSV file
	file, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// defining the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	// writing the column headers
	writer.Write(headers)

	// adding each Pokemon product to the CSV output file
	for _, v := range pokemonProducts {
		// converting a PokemonProduct to an array of strings
		record := []string{
			v.url,
			v.image,
			v.name,
			v.price,
		}
		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()

	return writer.Error()
}

func contains(pagesDiscovered []string, newPageLink string) bool {
	for _, page := range pagesDiscovered {
		if strings.EqualFold(page, newPageLink) {
			return true
		}
	}
	return false
}
