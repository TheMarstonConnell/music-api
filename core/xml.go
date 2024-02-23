package core

import (
	"encoding/xml"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

type Tag struct {
	Name string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}

type Conversions struct {
	Cube []Tag `xml:"Cube>Cube>Cube"`
}

func Convert(fromAmount float64, fromCurrency string, toCurrency string) float64 {
	m := getConv()
	f := m[strings.ToUpper(fromCurrency)]
	t := m[strings.ToUpper(toCurrency)]

	return (fromAmount / f) * t

}

func getConv() map[string]float64 {

	m := make(map[string]float64)

	if xmlBytes, err := getXML("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"); err != nil {
		log.Error().Err(fmt.Errorf("failed to get XML: %w", err))
	} else {

		var result Conversions
		err := xml.Unmarshal(xmlBytes, &result)
		if err != nil {
			log.Error().Err(err)
			return m
		}

		for _, s := range result.Cube {

			f, err := strconv.ParseFloat(s.Rate, 64)
			if err != nil {
				log.Error().Err(err)
				continue
			}
			m[s.Name] = f
		}

	}

	return m
}
