package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/teros0/cservices/ingestor"
	"github.com/teros0/cservices/storage"
)

var (
	port    = flag.String("port", ":7777", "http port to listen on")
	csvPath = flag.String("path", "./resources/data.csv", "path to the csv file")
)

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGKILL)
	fmt.Println("starting")
	go func() {
		select {
		case <-stop:
			cancelFunc()
		case <-ctx.Done():
		}
	}()

	wg.Add(2)

	go func() {
		defer func() {
			wg.Done()
		}()
		log.Println("running storage service")

		if err := (&storage.Storage{}).Run(ctx, *port); err != nil {
			log.Fatalf("while running storage service -> %s", err)
		}
	}()
	go func() {
		defer func() {
			wg.Done()
		}()
		log.Println("running ingestor service")

		if err := (&ingestor.Ingestor{}).Run(ctx, *csvPath, *port); err != nil {
			log.Fatalf("while running ingestor service -> %s", err)
		}
	}()
	wg.Wait()
	log.Println("finished")
}
