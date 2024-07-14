package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Please enter options")
		os.Exit(1)
	}

	fileName := os.Args[2]
	max, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println("error max")
		os.Exit(1)
	}
	cur := 0
	f, err := os.Create("./target/" + fileName)
	if err != nil {
		log.Println("Fail to create file")
		os.Exit(1)
	}
	w := bufio.NewWriter(f)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnHTML("#chapter-c", func(h *colly.HTMLElement) {
		// Replace <br> by new line
		h.DOM.Find("br").ReplaceWithHtml("<span>\n</span>")
		chapName := strings.Split(h.Request.URL.Path, "/")[2]
		_, err = w.WriteString(chapName + "\n" + h.DOM.Text() + "\n")
		if err != nil {
			log.Printf("Fail to write to file for %v", h.Request.URL.String())
			os.Exit(1)
		}
		w.Flush()
	})

	c.OnHTML("#next_chap", func(h *colly.HTMLElement) {
		if strings.Contains(h.Attr("class"), "disabled") {
			// No more chapter
			return
		}
		if cur > max {
			return
		}
		nextChapUrl := h.Attr("href")
		c.Visit(nextChapUrl)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Success: %v\n", r.Request.URL.String())
		cur += 1
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v\n", r.Request.URL.String())
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %v\n", r.URL.String())
	})

	c.Visit("https://truyenfull.vn/gia-thien/chuong-1/")
}
