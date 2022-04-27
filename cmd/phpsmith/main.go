package main

import (
	"fmt"
	"log"

	"github.com/cespare/subcmd"
)

// Build* variables are initialized during the build via -ldflags.
var (
	BuildVersion string
	BuildTime    string
	BuildOSUname string
	BuildCommit  string
)

func main() {
	log.SetFlags(0)

	cmds := []subcmd.Command{
		{
			Name:        "version",
			Description: "print phpsmith version info to stdout and exit",
			Do:          versionMain,
		},

		{
			Name:        "fuzz",
			Description: "run fuzzing using the provided configuration",
			Do:          fuzzMain,
		},

		{
			Name:        "generate",
			Description: "generate a program using the provided configuration",
			Do:          generateMain,
		},
	}

	subcmd.Run(cmds)
}

func versionMain(args []string) {
	if BuildCommit == "" {
		fmt.Printf("phpsmith built without version info\n")
	} else {
		fmt.Printf("phpsmith version %s\nbuilt on: %s\nos: %s\ncommit: %s\n",
			BuildVersion, BuildTime, BuildOSUname, BuildCommit)
	}
}

func fuzzMain(args []string) {
	if err := cmdFuzz(args); err != nil {
		log.Fatalf("phpsmith fuzz: error: %v", err)
	}
}

func generateMain(args []string) {
	if err := cmdGenerate(args); err != nil {
		log.Fatalf("phpsmith generate: error: %v", err)
	}
}
