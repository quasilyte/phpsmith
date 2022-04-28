package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

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

	interrupt := make(chan os.Signal, 1)
	signalNotify(interrupt)

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-interrupt
		cancel()
	}()

	filesCh := make(chan string, 100)
	for i := 0; i < concurrency; i++ {
		eg.Go(func() error {
			return runner(ctx, filesCh, seed)
		})
	}

out:
	for {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		files, err := generate(dir, seed)
		if err != nil {
			return err
		}

		for _, file := range files {
			select {
			case filesCh <- file:
			case <-ctx.Done():
				break out
			}
		}

		if err = eg.Wait(); err != nil {
			log.Println("on errorGroup Execution: ", err)
		}
	}

	return nil
}

func runner(ctx context.Context, files <-chan string, seed int64) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case filename := <-files:
			if err := fuzzingProcess(ctx, filename, seed); err != nil {
				log.Println("on fuzzingProcess:", err)
			}
		}
	}
}

type ExecutorOutput struct {
	Output string
	Error  string
}

func fuzzingProcess(ctx context.Context, filename string, seed int64) error {
	var results = make(map[int]ExecutorOutput, len(executors))
	var errMsg string
	for i, ex := range executors {
		r, err := ex(ctx, filename)

		if err != nil {
			errMsg = err.Error()
		}

		results[i] = ExecutorOutput{
			Output: string(r),
			Error:  errMsg,
		}
	}

	checkedStack := make(map[int][]int)
	for i, r := range results {
	inner:
		for ii, rr := range results {
			if checks, ok := checkedStack[i]; ok {
				for _, checkedRes := range checks {
					if checkedRes == ii {
						continue inner
					}
				}
			}

			if diff := cmp.Diff(r.Output, rr.Output); diff != "" {
				checkedStack[i] = append(checkedStack[i], ii)
				log.Println("-----------------------------")
				log.Printf("out: %s\tseed: %d\tdiff: %s\tstdErr: %s\t \n", r.Output, seed, diff, r.Error)
				log.Printf("out: %s\tseed: %d\tdiff: %s\tstdErr: %s\t \n", rr.Output, seed, diff, rr.Error)
			}
		}
	}

	return nil
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}
