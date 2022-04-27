package main

import "flag"

func cmdFuzz(args []string) error {
	fs := flag.NewFlagSet("phpsmith fuzz", flag.ExitOnError)
	_ = fs.Parse(args)

	// TODO

	return nil
}
