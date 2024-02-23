package core

import (
	"github.com/rs/zerolog/log"
	"sort"
)

type Store struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Link  string  `json:"link"`
}

type Response struct {
	Stores []Store `json:"stores"`
}

func NewResponse() *Response {
	r := Response{make([]Store, 0)}
	return &r

}

func (r *Response) AddStore(name string, link string, price float64) {
	r.Stores = append(r.Stores, Store{
		Name:  name,
		Price: price,
		Link:  link,
	})

	sort.Slice(r.Stores, func(i, j int) bool {
		return r.Stores[i].Price < r.Stores[j].Price
	})
}

func GetPrice(search string) *Response {
	r := NewResponse()

	bcprice, bcl, err := BCSearchAlbum(search, "cad")
	if err != nil || len(bcl) == 0 {
		log.Warn().Err(err)
	} else {
		r.AddStore("Bandcamp", bcl, bcprice)
	}

	quoprice, ql, err := QuoSearchAlbum(search, "ca")
	if err != nil || len(ql) == 0 {
		log.Warn().Err(err)
	} else {
		r.AddStore("Quobuz", ql, quoprice)
	}

	return r
}
