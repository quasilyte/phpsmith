package main

import (
	"bytes"
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

type executor func(ctx context.Context, filename string, seed int64) ([]byte, error)

var executors = []executor{
	interpretator.RunPHP,
	interpretator.RunKPHP,
}

func cmdFuzz(args []string) error {
	fs := flag.NewFlagSet("phpsmith fuzz", flag.ExitOnError)
	flagConcurrency := fs.Int("concurrency", 1,
		"Number of concurrent runners. Defaults to the half number of available CPU cores.")
	flagOutputDir := fs.String("o", "phpsmith_out",
		`output dir`)

	_ = fs.Parse(args)

	concurrency := *flagConcurrency
	dir := *flagOutputDir

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

	dirCh := make(chan dirAndSeed, concurrency)
	for i := 0; i < concurrency; i++ {
		eg.Go(func() error {
			return runner(ctx, dirCh)
		})
	}

	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
out:
	for {
		seed := randomizer.Int63()
		newDir := dir + "_" + strconv.FormatInt(seed, 10)
		if err := generate(newDir, seed); err != nil {
			log.Println("on generate: ", err)
			continue
		}

		select {
		case dirCh <- dirAndSeed{Dir: newDir, Seed: seed}:
		case <-ctx.Done():
			break out
		}
	}

	if err := eg.Wait(); err != nil {
		log.Println("on errorGroup Execution: ", err)
	}

	return nil
}

func runner(ctx context.Context, dirCh <-chan dirAndSeed) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case ds := <-dirCh:
			diffFound := fuzzingProcess(ctx, ds)
			suffix := ""
			if diffFound {
				suffix = "(found diff)"
			} else {
				if err := os.RemoveAll(ds.Dir); err != nil {
					return err
				}
			}
			log.Println("dir processed:", ds.Dir, suffix)
		}
	}
}

type executorOutput struct {
	Output string
	Error  string
}

type dirAndSeed struct {
	Dir  string
	Seed int64
}

func fuzzingProcess(ctx context.Context, ds dirAndSeed) bool {
	var results = make(map[int]executorOutput, len(executors))
	for i, ex := range executors {
		var errMsg string
		r, err := ex(ctx, ds.Dir, ds.Seed)

		if err != nil {
			errMsg = err.Error()
		}

		grepExceptions(r, ds.Seed)
		results[i] = executorOutput{
			Output: string(r),
			Error:  errMsg,
		}
	}

	diff := cmp.Diff(results[0].Output, results[1].Output)
	if diff != "" {
		l, err := os.OpenFile("./"+ds.Dir+"/log", os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			log.Println("-----------------------------")
			log.Printf("diff: %s\t, seed: %d\t\n", diff, ds.Seed)
			log.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
			log.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
		} else {
			defer l.Close()

			logger := log.New(l, "", 0)
			logger.Printf("diff: %s\t, seed: %d\t\n", diff, ds.Seed)
			logger.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
			logger.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
		}
	}
	return diff != ""
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}

var exceptionPatterns = [][]byte{
	[]byte("uncaught exception"),
	[]byte("fatal error"),
}

func grepExceptions(s []byte, seed int64) {
	for _, pattern := range exceptionPatterns {
		if bytes.Contains(bytes.ToLower(s), pattern) {
			log.Println("found exception pattern: on seed:", seed)
		}
	}
}
