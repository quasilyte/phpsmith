package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"

	"github.com/quasilyte/phpsmith/cmd/phpsmith/interpretator"
)

type executor func(ctx context.Context, filename string) ([]byte, error)

var executors = []executor{
	interpretator.RunPHP,
	interpretator.RunKPHP,
}

func cmdFuzz(args []string) error {
	fs := flag.NewFlagSet("phpsmith fuzz", flag.ExitOnError)
	flagConcurrency := fs.Int("flagConcurrency", 1,
		"Number of concurrent runners. Defaults to the half number of available CPU cores.")
	flagSeed := fs.Int64("seed", 0,
		`a seed to be used during the code generation, 0 means "randomized seed"`)
	flagOutputDir := fs.String("o", "phpsmith_out",
		`output dir`)

	_ = fs.Parse(args)

	concurrency := *flagConcurrency
	dir := *flagOutputDir
	seed := *flagSeed

	if concurrency == 1 && runtime.NumCPU()/2 > 1 {
		concurrency = runtime.NumCPU() / 2
	}

	if seed == 0 {
		seed = time.Now().Unix()
	}

	interrupt := make(chan os.Signal, 1)
	signalNotify(interrupt)

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-interrupt
		cancel()
	}()

	dirCh := make(chan string, 100)
	for i := 0; i < concurrency; i++ {
		eg.Go(func() error {
			return runner(ctx, dirCh, seed)
		})
	}

	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
out:
	for {
		newDir := dir + strconv.Itoa(randomizer.Int())
		if err := generate(newDir, seed); err != nil {
			log.Println("on generate: ", err)
			continue
		}

		select {
		case dirCh <- newDir:
		case <-ctx.Done():
			break out
		}
	}

	if err := eg.Wait(); err != nil {
		log.Println("on errorGroup Execution: ", err)
	}

	return nil
}

func runner(ctx context.Context, dirs <-chan string, seed int64) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case dir := <-dirs:
			if err := fuzzingProcess(ctx, dir, seed); err != nil {
				log.Println("on fuzzingProcess:", err)
			}
			log.Println("dir processed:", dir)
		}
	}
}

type ExecutorOutput struct {
	Output string
	Error  string
}

func fuzzingProcess(ctx context.Context, dir string, seed int64) error {
	var results = make(map[int]ExecutorOutput, len(executors))
	var errMsg string
	for i, ex := range executors {
		r, err := ex(ctx, dir)

		if err != nil {
			errMsg = err.Error()
		}

		results[i] = ExecutorOutput{
			Output: string(r),
			Error:  errMsg,
		}
	}

	if diff := cmp.Diff(results[0].Output, results[1].Output); diff != "" {
		l, err := os.OpenFile(dir+"/log", os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			log.Println("-----------------------------")
			log.Printf("diff: %s\t, seed: %d\t\n", diff, seed)
			log.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
			log.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
		} else {
			defer l.Close()

			logger := log.New(l, "", 0)
			logger.Printf("diff: %s\t, seed: %d\t\n", diff, seed)
			logger.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
			logger.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
		}
	}

	return nil
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}
