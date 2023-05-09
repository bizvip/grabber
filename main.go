package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"strconv"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Hour)
	for {
		crawlNews()
		<-ticker.C
	}
}

func crawlNews() {
	c := colly.NewCollector()

	c.OnHTML("li", func(e *colly.HTMLElement) {
		category := e.ChildText("div.dd_lm")
		title := e.ChildText("div.dd_bt")
		timeStr := e.ChildText("div.dd_time")

		if category != "" && title != "" && timeStr != "" {
			fmt.Printf("Category: %s, Title: %s, Time: %s\n", category, title, timeStr)
		} else {
			return
		}

		t, err := time.Parse("1-2 15:04", timeStr)
		if err != nil {
			fmt.Println("时间无法解析:", err)
			return
		}

		now := time.Now()
		year := now.Year()

		t = time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		formattedTime := t.Format("20060102")

		fileName := "/www/wwwroot/156.241.140.23/text/xwbt/" + formattedTime + ".txt"

		// Check if title already exists in file
		if titleExistsInFile(fileName, title) {
			return
		}

		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("无法打开要写入的txt文件:", err)
			return
		}
		defer f.Close()

		encoder := simplifiedchinese.GBK.NewEncoder()
		encodedTitle, _ := encoder.String(title)

		if _, err := f.WriteString(encodedTitle + "\n"); err != nil {
			fmt.Println("无法写入并转码文件:", err)
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	baseURL := "https://www.chinanews.com/scroll-news/news"
	for i := 1; i <= 10; i++ {
		_ = c.Visit(baseURL + strconv.Itoa(i) + ".html")
		time.Sleep(5 * time.Second)
	}
}

func titleExistsInFile(fileName, title string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer file.Close()

	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(file, decoder)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if scanner.Text() == title {
			return true
		}
	}

	return false
}
