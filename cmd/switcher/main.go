package main

import (
	"os"
	"os/signal"
	"switcher/internal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// pred init logger
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    true,
		TimeFormat: time.DateTime,
	})

	// Working directory for nix - pwd()
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	conf := internal.LoadConfig(workDir)

	// set global log level from config
	if conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	app := internal.NewApp(conf)
	graceChan := make(chan os.Signal, 1)
	signal.Notify(graceChan, os.Interrupt)
	go func() {
		_ = <-graceChan
		_ = app.Shutdown()
	}()

	_ = internal.WithLogFile(workDir, true, func() error {
		internal.PrintVersion()

		if err := app.Serve(); err != nil {
			log.Error().Msgf("%v", err)
		}
		return err
	})
}
