package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

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
		flagConcurrency = ptrOfInt(runtime.NumCPU() / 2)
	}

	interrupt := make(chan os.Signal)
	signalNotify(interrupt)

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-interrupt
		cancel()
	}()

	for i := 0; i < *flagConcurrency; i++ {
		eg.Go(func() error {
			return runner(ctx, dir, seed)
		})
	}

	return eg.Wait()
}

func runner(ctx context.Context, dir string, seed int64) error {
	for {
		files, err := generate(dir, seed)
		if err != nil {
			return err
		}

		for _, filename := range files {
			if err = fuzzingProcess(ctx, filename); err != nil {
				return err
			}
		}
	}
}

func fuzzingProcess(ctx context.Context, filename string) error {
	results := make([][]byte, 0, len(executors))
	errors := make([]error, 0, len(executors))

	for _, ex := range executors {
		r, err := ex(ctx, filename)

		results = append(results, r)
		errors = append(errors, err)
	}

	return compareResults(results, errors)
}

func compareResults(res [][]byte, errors []error) error {
	return nil
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func ptrOfInt(i int) *int {
	return &i
}
