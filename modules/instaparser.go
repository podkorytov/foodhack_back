package modules

import (
	"github.com/PuerkitoBio/goquery"
	"errors"
)

func GetInstaImage(link string) (string, error)  {
	doc, _ := goquery.NewDocument(link)
	item := ""

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		op, _ := s.Attr("property")
		url, _ := s.Attr("content")

		if op == "og:image" {
			item = url
		}
	})

	if item != "" {
		return  item, nil
	}

	return "", errors.New("can't parse url")
}