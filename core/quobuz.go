package core

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

func QuoGetPrice(u *url.URL) (float64, error) {
	class := ".album-addtocart__highlight"

	c := colly.NewCollector()

	priceVal := 0.0

	// Find and visit all links
	c.OnHTML(class, func(e *colly.HTMLElement) {
		price := e.Text

		s := strings.Replace(price, "$", "", -1)
		s = strings.Replace(s, "CA", "", -1)
		s = strings.TrimSpace(s)

		p, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return
		}
		priceVal = p
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debug().Msgf("Visiting %s", r.URL)
	})

	err := c.Visit(u.String())
	if err != nil {
		return 0.0, err
	}
	return priceVal, nil
}

func QuoSearchAlbum(album string, region string) (float64, string, error) {
	// path := "https://www.qobuz.com/ca-en/search?q=janelle+monae+archandroid"

	base := "https://www.qobuz.com"
	b, err := url.Parse(base)
	if err != nil {
		return 0.0, "", fmt.Errorf("cannot find album %w", err)
	}

	ablm := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(album, "")

	u := b.JoinPath(fmt.Sprintf("%s-en", region), "search")
	q := u.Query()
	q.Set("q", ablm)
	q.Set("s", "prc")
	u.RawQuery = q.Encode()

	c := colly.NewCollector()

	price := 0.0
	l := ""

	found := false

	// Find and visit all links
	c.OnHTML(".btn__qobuz--see-album", func(e *colly.HTMLElement) {
		if found {
			return
		}
		link := e.Attr("href")
		newLink := b.JoinPath(link)
		log.Debug().Msgf("searching at %s", newLink.String())
		p, err := QuoGetPrice(newLink)
		if err != nil {
			log.Error().Err(err)
			return
		}
		price = p
		found = true
		l = newLink.String()
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debug().Msgf("Visiting %s", r.URL)
	})

	err = c.Visit(u.String())
	if err != nil {
		return price, l, err
	}
	return price, l, nil
}
