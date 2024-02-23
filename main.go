package main

import (
	"os"

	"github.com/TheMarstonConnell/music-api/net"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Level(zerolog.DebugLevel)

	net.Start()
}
