package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func Execute(fileName, novelName string, total, batchSize int) {
	var wg sync.WaitGroup
	nWorker := total / batchSize
	if total%batchSize > 0 {
		nWorker++
	}
	dir := fmt.Sprintf("./target/%v", novelName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			os.Exit(1)
		}
	}
	for i := 0; i < nWorker; i++ {
		wg.Add(1)
		go func() {
			start := i*batchSize + 1
			end := start + batchSize - 1
			workerName := fmt.Sprintf("Crawler %d-%d", start, end)
			hasError := false
			if end > total {
				end = total
			}
			file, err := os.Create(fmt.Sprintf("%v/%v_%d-%d.txt", dir, fileName, start, end))
			if err != nil {
				log.Printf("[%v] Fail to create file %v\n", workerName, err)
				wg.Done()
				return
			}
			w := bufio.NewWriter(file)

			c := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
			)

			c.OnHTML("#chapter-c", func(h *colly.HTMLElement) {
				// Replace <br> by new line
				h.DOM.Find("br").ReplaceWithHtml("<span>\n</span>")
				chapName := strings.Split(h.Request.URL.Path, "/")[2]
				_, err = w.WriteString(chapName + "\n" + h.DOM.Text() + "\n")
				if err != nil {
					log.Printf("[%v] Fail to write to file for %v", workerName, h.Request.URL.String())
				}
				w.Flush()
			})

			c.OnHTML("#next_chap", func(h *colly.HTMLElement) {
				if strings.Contains(h.Attr("class"), "disabled") {
					// No more chapter
					return
				}
				if start > end {
					return
				}
				nextChapUrl := h.Attr("href")
				c.Visit(nextChapUrl)
			})

			c.OnResponse(func(r *colly.Response) {
				start += 1
			})

			c.OnError(func(r *colly.Response, err error) {
				log.Printf("[%v] Error: %v\n%v", workerName, r.Request.URL.String(), err.Error())
				hasError = true
			})

			c.Visit(fmt.Sprintf("https://truyenfull.io/%v/chuong-%d/", novelName, start))
			if !hasError {
				log.Printf("[%v] Success\n", workerName)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
