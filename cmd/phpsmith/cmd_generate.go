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
	flagSeed := fs.Int64("seed", 0,
		`a seed to be used during the code generation, 0 means "randomized seed"`)
	flagOutputDir := fs.String("o", "phpsmith_out",
		`output dir`)
	_ = fs.Parse(args)

	_, err := generate(*flagOutputDir, *flagSeed)

	return err
}

func generate(dir string, seed int64) ([]string, error) {
	randomSeed := seed
	if randomSeed == 0 {
		randomSeed = time.Now().Unix()
	}
	random := rand.New(rand.NewSource(randomSeed))

	if err := os.MkdirAll(dir, 0o700); err != nil && !os.IsExist(err) {
		return nil, err
	}

	config := &irgen.Config{Rand: random}
	program := irgen.CreateProgram(config)
	printerConfig := &irprint.Config{
		Rand: random,
	}

	var filenames []string
	for _, f := range program.Files {
		fullname := filepath.Join(dir, f.Name)
		filenames = append(filenames, fullname)

		fileContents := makeFileContents(f, printerConfig)
		if err := os.WriteFile(fullname, fileContents, 0o664); err != nil {
			return nil, fmt.Errorf("create %s file: %w", fullname, err)
		}
	}

	return filenames, nil
}

func makeFileContents(f *irgen.File, config *irprint.Config) []byte {
	var buf bytes.Buffer
	buf.WriteString("<?php\n")
	for _, n := range f.Nodes {
		irprint.FprintRootNode(&buf, n, config)
	}
	return buf.Bytes()
}
