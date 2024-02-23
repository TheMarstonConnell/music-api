package core

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"strings"
)

func BCGetPrice(u *url.URL, inCurrency string) (float64, error) {
	class := ".digital > .ft > .ft > .nobreak"

	c := colly.NewCollector()

	var priceVal = 0.0
	complete := false
	foundAnything := false

	// Find and visit all links
	c.OnHTML(class, func(e *colly.HTMLElement) {
		foundAnything = true
		price := e.ChildText(".base-text-color")
		currency := e.ChildText(".buyItemExtra")
		currency = currency[:len(currency)-9]

		s := strings.Replace(price, "$", "", -1)
		p, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Error().Err(err)
			return
		}

		priceVal = Convert(p, currency, inCurrency)
		complete = true
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debug().Msgf("Visiting %s", r.URL)
	})

	err := c.Visit(u.String())
	if err != nil {
		return 0.0, fmt.Errorf("failed to get price: %w", err)
	}

	if !complete || !foundAnything {
		return priceVal, fmt.Errorf("could not get price")
	}

	return priceVal, nil
}

func BCSearchAlbum(album string, currency string) (float64, string, error) {
	base := "https://bandcamp.com"
	b, err := url.Parse(base)
	if err != nil {
		return 0.0, "", fmt.Errorf("cannot find album %w", err)
	}

	u := b.JoinPath("search")
	q := u.Query()
	q.Set("from", "autocomplete")
	q.Set("q", album)
	q.Set("item_type", "a")
	u.RawQuery = q.Encode()

	c := colly.NewCollector()

	found := false
	bcl := ""
	var p = 0.0

	// Find and visit all links
	c.OnHTML(".result-info > .heading > a", func(e *colly.HTMLElement) {
		if found {
			return
		}
		link := e.Attr("href")

		l, err := url.Parse(link)
		if err != nil {
			return
		}
		price, err := BCGetPrice(l, currency)
		if err != nil {
			return
		}
		p = price
		found = true
		bcl = l.String()
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debug().Msgf("Visiting %s", r.URL)
	})

	c.Visit(u.String())

	return p, bcl, nil
}
