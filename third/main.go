package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	siteURL = "https://hypeauditor.com/top-instagram-all-russia/"
)

// parseSite будет парсить указанный сайт и возвращать данные в виде слайса слайсов.
func parseSite() ([][]string, error) {

	resp, err := http.Get(siteURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var data [][]string
	doc.Find("div.row__top").Each(func(i int, row *goquery.Selection) {
		rank := row.Find("div.row-cell.rank span").First().Text()
		nick := row.Find("div.contributor-wrap a.contributor div.contributor__name-content").Text()
		name := row.Find("div.contributor-wrap a.contributor div.contributor__title").Text()
		categories := row.Find("div.row-cell.category div.tag__content").Map(func(_ int, sel *goquery.Selection) string {
			return sel.Text()
		})
		subscribers := row.Find("div.row-cell.subscribers").Text()

		data = append(data, []string{
			strings.TrimSpace(rank),
			strings.TrimSpace(nick),
			strings.TrimSpace(name),
			strings.Join(categories, ", "),
			strings.TrimSpace(subscribers),
		})
	})

	return data, nil
}

func main() {
	data, err := parseSite()
	if err != nil {
		log.Fatalf("Error parsing the site: %s", err)
		return
	}

	file, err := os.Create("sample_result.csv")
	if err != nil {
		log.Fatalf("Error creating the CSV file: %s", err)
		return
	}
	defer file.Close()

	// Разделитель - ';'
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	writer.Write([]string{"Рейтинг", "Ник", "Имя", "Категория", "Подписчики"})

	for _, row := range data {
		writer.Write(row)
	}

	fmt.Println("Parsing completed and data saved to sample_result.csv")
}
