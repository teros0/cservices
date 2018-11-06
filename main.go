package main

import (
	"context"
	"cservices/ingestor"
	"cservices/storage"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	port    = flag.String("port", ":7777", "http port to listen on")
	csvPath = flag.String("path", "./resorces/data.csv", "path to the csv file")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGKILL)

	go func() {
		select {
		case <-stop:
			cancelFunc()
		case <-ctx.Done():
		}
	}()

	go func() {
		log.Println("running storage service")

		if err := (&storage.Storage{}).Run(ctx, *port); err != nil {
			log.Fatalf("while running storage service -> %s", err)
		}

		log.Println("running ingestor service")

		if err := (&ingestor.Ingestor{}).Run(ctx, *csvPath, *port); err != nil {
			log.Fatalf("while running ingestor service -> %s", err)
		}
	}()
}
