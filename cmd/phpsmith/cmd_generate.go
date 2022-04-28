package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/quasilyte/phpsmith/irgen"
	"github.com/quasilyte/phpsmith/irprint"
)

func cmdGenerate(args []string) error {
	fs := flag.NewFlagSet("phpsmith generate", flag.ExitOnError)
	flagSeed := fs.Uint64("seed", 0,
		`a seed to be used during the code generation, 0 means "randomized seed"`)
	flagOutputDir := fs.String("o", "phpsmith_out",
		`output dir`)
	_ = fs.Parse(args)

	randomSeed := int64(*flagSeed)
	if randomSeed == 0 {
		randomSeed = time.Now().Unix()
	}
	random := rand.New(rand.NewSource(randomSeed))

	if err := os.MkdirAll(*flagOutputDir, 0o700); err != nil {
		return err
	}

	config := &irgen.Config{Rand: random}
	program := irgen.CreateProgram(config)
	printerConfig := &irprint.Config{
		Rand: random,
	}
	for _, f := range program.RuntimeFiles {
		fullname := filepath.Join(*flagOutputDir, f.Name)
		if err := os.WriteFile(fullname, f.Contents, 0o664); err != nil {
			return fmt.Errorf("create %s file: %w", fullname, err)
		}
	}
	for _, f := range program.Files {
		fullname := filepath.Join(*flagOutputDir, f.Name)
		fileContents := makeFileContents(f, printerConfig)
		if err := os.WriteFile(fullname, fileContents, 0o664); err != nil {
			return fmt.Errorf("create %s file: %w", fullname, err)
		}
	}

	return nil
}

func makeFileContents(f *irgen.File, config *irprint.Config) []byte {
	var buf bytes.Buffer
	buf.WriteString("<?php\n")
	for _, n := range f.Nodes {
		irprint.FprintRootNode(&buf, n, config)
	}
	return buf.Bytes()
}
