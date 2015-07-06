package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

func getBalance(doc *goquery.Document) (int, error) {
	selector := fmt.Sprintf("img[src=\"%s\"]", balanceImageSrc)
	balanceStr := doc.Find(selector).First().Siblings().Text()
	cleanStr := regexp.MustCompile("-?\\d+").FindString(regexp.MustCompile("\\s*").ReplaceAllString(balanceStr, ""))
	if len(cleanStr) == 0 {
		return 0, fmt.Errorf(
			"Cannot get byfly data. Check login and password, try again within 10 minutes.")
	}
	return strconv.Atoi(cleanStr)
}

func getTariff(doc *goquery.Document) string {
	selector := fmt.Sprintf("td:contains(\"%s\")", tariffText)
	return doc.Find(selector).First().Siblings().First().Text()
}

func getClient(doc *goquery.Document) string {
	selector := fmt.Sprintf("td:contains(\"%s\")", clientText)
	return doc.Find(selector).First().Siblings().First().Text()
}

func getStatus(doc *goquery.Document) string {
	selector := fmt.Sprintf("td:contains(\"%s\")", statusText)
	return doc.Find(selector).First().Siblings().First().Text()
}

func (byfly *Byfly) parsePage(body string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}
	// balance
	balance, err := getBalance(doc)
	if err != nil {
		return err
	}
	byfly.balance = balance
	byfly.tariff = getTariff(doc)
	byfly.client = getClient(doc)
	byfly.status = getStatus(doc)
	return nil
}
