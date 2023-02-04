package main

import (
	"context"
	"elastic-tools/internal/elastic"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var action string
	var host string
	var index string
	var periodSeconds int
	var batch int
	var size int

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "action",
				Usage:       "info, insert",
				Value:       "info",
				Destination: &action,
			},
			&cli.StringFlag{
				Name:        "host",
				Value:       "http://localhost:9200",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "index",
				Destination: &index,
			},
			&cli.IntFlag{
				Name:        "batch",
				Value:       1,
				Destination: &batch,
			},
			&cli.IntFlag{
				Name:        "period",
				Value:       1,
				Destination: &periodSeconds,
			},
			&cli.IntFlag{
				Name:        "size",
				Value:       1,
				Destination: &size,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Printf("action=%v, host=%v, index=%v, periodSeconds=%v\n", action, host, index, periodSeconds)
			elasticsearch, err := elastic.CreateClient([]string{host})
			if err != nil {
				os.Exit(1)
			}
			switch action {
			case "info":
				info(elasticsearch)
			case "insert":
				insertDocument(elasticsearch, index, batch, size, periodSeconds)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func info(es *elastic.Elastic) {
	elastic.Info(es)
}

func insertDocument(es *elastic.Elastic, index string, batch, size, sleep int) error {
	if index == "" {
		log.Println("index is not valid")
		return nil
	}
	ctx := context.Background()
	elastic.InsertDummyDocument(es, ctx, index, batch, size, sleep)
	return nil
}
