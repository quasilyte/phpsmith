package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/google/go-cmp/cmp"
	"github.com/quasilyte/phpsmith/cmd/phpsmith/interpretator/kphp"
	"github.com/quasilyte/phpsmith/cmd/phpsmith/interpretator/php"
)

type Runner interface {
	Run(ctx context.Context, filename string, seed int64) ([]byte, error)
	Name() string
}

var runners = []Runner{
	php.Runner{},
	kphp.Runner{},
}

func cmdFuzz(args []string) error {
	fs := flag.NewFlagSet("phpsmith fuzz", flag.ExitOnError)
	flagConcurrency := fs.Int("concurrency", 0,
		"Number of concurrent runners. Defaults to the half number of available CPU cores.")
	flagOutputDir := fs.String("o", "phpsmith_out",
		`output dir`)

	_ = fs.Parse(args)

	concurrency := *flagConcurrency
	dir := *flagOutputDir

	if concurrency == 0 {
		concurrency = 1
		if runtime.NumCPU()/2 > 1 {
			concurrency = runtime.NumCPU() / 2
		}
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
		return fmt.Errorf("on errorGroup Execution: %w", err)
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
	var (
		results = make([]executorOutput, 0, len(runners))
		wg      sync.WaitGroup
	)

	for _, r := range runners {
		var (
			err   error
			out   []byte
			errCh = make(chan error)
		)

		wg.Add(1)

		func() { // anon func need for close context
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				defer wg.Done()
				out, err = r.Run(ctx, ds.Dir, ds.Seed)
				errCh <- err
			}()

			select {
			case err = <-errCh:
			case <-time.After(time.Minute):
				err = fmt.Errorf("too long execution for: %s on seed %d", r.Name(), ds.Seed)
				cancel()
			}
		}()

		wg.Wait()

		grepExceptions(out, ds.Seed)

		var msg string
		if err != nil {
			msg = err.Error()
		}

		results = append(results, executorOutput{
			Output: string(out),
			Error:  msg,
		})
	}

	writeLog := func(diff string) {
		l, err := os.OpenFile("./"+ds.Dir+"/log", os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			log.Println("-----------------------------")
			log.Printf("diff: %s\t, seed: %d\t\n", diff, ds.Seed)
			log.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
			log.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
			return
		}
		defer l.Close()

		logger := log.New(l, "", 0)
		logger.Printf("diff: %s\t, seed: %d\t\n", diff, ds.Seed)
		logger.Printf("out: %s\terr: %s\t\n", results[0].Output, results[0].Error)
		logger.Printf("out: %s\terr: %s\t\n", results[1].Output, results[1].Error)
	}

	diff := cmp.Diff(results[0].Output, results[1].Output)
	if diff != "" {
		writeLog(diff)
	} else if results[0].Error != "" || results[1].Error != "" {
		diff = "nil"
		writeLog(diff)
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
