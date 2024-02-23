package net

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/TheMarstonConnell/music-api/core"
	"github.com/TheMarstonConnell/music-api/net/templates"
	"github.com/a-h/templ"
	"github.com/rs/zerolog/log"
)

//go:embed public/*
var public embed.FS

func GetAlbum() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		res := core.GetPrice(q)

		s := templates.StoreList(res)
		_ = s.Render(context.Background(), w)
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(res)
	}
}

func Start() {
	pub, err := fs.Sub(public, "public")
	if err != nil {
		log.Error().Err(err)
		return
	}

	entries, err := public.ReadDir("public")
	if err != nil {
		log.Error().Err(err)
		return
	}

	css := make([]string, 0)

	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".css") {
			css = append(css, fmt.Sprintf("/public/%s", entry.Name()))
		}
	}

	component := templates.App(css)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/search", GetAlbum())

	fs := http.FileServer(http.FS(pub))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	mux.Handle("/", templ.Handler(component))

	s := &http.Server{
		Addr:           ":9797",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Error().Err(err)
	}
}
