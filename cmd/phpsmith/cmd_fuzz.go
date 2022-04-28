package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func cmdFuzz(args []string) error {
	fs := flag.NewFlagSet("phpsmith fuzz", flag.ExitOnError)
	concurrency := fs.Int("concurrency", 1, "Number of concurrent runners. Defaults to the half number of available CPU cores.")
	_ = fs.Parse(args)

	if *concurrency == 1 && runtime.NumCPU()/2 > 1 {
		concurrency = ptrOfInt(runtime.NumCPU() / 2)
	}

	var (
		wg            = &sync.WaitGroup{}
		interrupt     = make(chan os.Signal)
		routinesLimit = make(chan struct{}, *concurrency)
	)
	signalNotify(interrupt)

	ctx, cancel := context.WithCancel(context.Background())

out:
	for {
		select {
		case <-interrupt:
			cancel()
			break out
		case routinesLimit <- struct{}{}:
		}

		wg.Add(1)
		go worker(ctx, routinesLimit, wg, fuzzingProcess)
	}

	wg.Wait()

	return nil
}

func worker(ctx context.Context, routinesLimit chan struct{}, wg *sync.WaitGroup, f func(ctx context.Context) error) {
	defer func() {
		wg.Done()
		<-routinesLimit
	}()

	if err := f(ctx); err != nil {
		log.Println("on worker: ", err)
	}
}

func fuzzingProcess(ctx context.Context) error {
	return nil
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func ptrOfInt(i int) *int {
	return &i
}
