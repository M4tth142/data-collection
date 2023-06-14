package webscraping

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type Article struct {
	Title string
	Text  string
}

func Web() {
	url := "https://www.beursduivel.be/Nieuws/Default.aspx"
	articles := []Article{}

	c := colly.NewCollector(
		colly.AllowedDomains("beursduivel.be"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.StatusCode)
	})

	c.OnHTML("article#articlecontainer", func(e *colly.HTMLElement) {
		// Find the <h3> and <p> tags within the matched element
		e.ForEach("h3, p", func(index int, elem *colly.HTMLElement) {
			// Check if the current element is an <h3> tag
			if elem.Name == "h3" {
				// Create a new Article struct
				article := Article{Title: elem.Text}
				// Check if the next element is a <p> tag
				if index+1 < len(e.DOM.Nodes) && e.DOM.Nodes[index+1].Data == "p" {
					// Set the Text field of the Article struct
					article.Text = e.DOM.Nodes[index+1].FirstChild.Data
				}
				// Append the Article struct to the articles array
				articles = append(articles, article)
			}
		})
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Set cookies
	c.SetCookies(url, []*http.Cookie{
		{Name: "cookie_name", Value: "cookie_value"},
		// Add more cookies here if necessary
	})

	// Start the scraping process
	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// Print the scraped articles
	for _, article := range articles {
		fmt.Printf("Title: %s\nText: %s\n\n", article.Title, article.Text)
	}
}
