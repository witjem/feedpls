package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/go-pkgz/lgr"

	"github.com/urfave/cli/v2"

	"github.com/witjem/feedpls/internal/server"
)

var revision string

func main() {
	fmt.Printf("[INFO] feedpls %s\n", revision)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
		<-ch

		cancel()

		log.Printf("[INFO] gracefully shutdown")
	}()

	cfg := server.Config{Version: revision}
	cliApp := &cli.App{
		Name:        "nfp",
		Description: "News feed please",
		Usage:       "run http server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				EnvVars:     []string{"APP_PORT"},
				Value:       8080,
				Usage:       "server port",
				Destination: &cfg.Listen,
			},
			&cli.StringFlag{
				Name:        "secret",
				EnvVars:     []string{"APP_SECRET"},
				Value:       "secret",
				Usage:       "api secret key",
				Destination: &cfg.Secret,
			},
			&cli.DurationFlag{
				Name:        "ttl",
				EnvVars:     []string{"APP_TTL"},
				Value:       time.Minute * 5,
				Usage:       "feeds caching time",
				Destination: &cfg.CacheTTL,
			},
			&cli.PathFlag{
				Name:        "feeds",
				EnvVars:     []string{"APP_FEEDS"},
				Value:       "./feeds.yml",
				Usage:       "yaml file with describes feeds matching",
				Destination: &cfg.FeedsCfgFilePath,
			},
		},
		Action: func(c *cli.Context) error {
			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}

			return srv.Run(ctx)
		},
	}

	err := cliApp.RunContext(ctx, os.Args)
	if err != nil {
		log.Printf("[ERROR] failed to run server, %v", err)
	}
}
