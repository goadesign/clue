package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"goa.design/clue/log"
	cluemockgen "goa.design/clue/mock/cmd/cmg/pkg"
)

func main() {
	var (
		gSet           = flag.NewFlagSet("global", flag.ExitOnError)
		debug, help, h *bool
		addGlobals     = func(set *flag.FlagSet) {
			debug = set.Bool("debug", false, "Print debug output")
			help = set.Bool("help", false, "Print help information")
			h = set.Bool("h", false, "Print help information")
		}

		genSet = flag.NewFlagSet("gen", flag.ExitOnError)

		versionSet = flag.NewFlagSet("version", flag.ExitOnError)

		showUsage = func(code int) {
			printUsage(gSet)
			os.Exit(code)
		}
	)

	addGlobals(gSet)

	if len(os.Args) == 1 {
		showUsage(1)
	}

	var (
		cmd  = os.Args[1]
		args []string
	)
	switch cmd {
	case "gen":
		addGlobals(genSet)
		_ = genSet.Parse(os.Args[2:])
		args = genSet.Args()
	case "version":
		addGlobals(versionSet)
		_ = versionSet.Parse(os.Args[2:])
		args = versionSet.Args()
	case "help":
		showUsage(0)
	default:
		_ = gSet.Parse(os.Args[1:])
	}

	if *h || *help {
		showUsage(0)
	}

	switch cmd {
	case "gen":
		ctx := context.Background()
		if *debug {
			ctx = log.Context(ctx, log.WithDebug())
		} else {
			ctx = log.Context(ctx)
		}
		err := cluemockgen.Generate(ctx, args, "")
		if err != nil {
			os.Exit(1)
		}
	case "version":
		fmt.Println(os.Args[0], "version", cluemockgen.Version())
	default:
		fmt.Fprintf(os.Stderr, `unknown command %q, use "--help" for usage`, cmd)
		os.Exit(1)
	}
}

func printUsage(fss ...*flag.FlagSet) {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, `%v is the Clue Mock Generation tool for the Goa framework.

Usage:
  %v gen PACKAGE...
  %v version

Commands:
  gen
        Generate mocks for interfaces in packages
  version
        Print version information

Args:
  PACKAGE
        Go import path(s) to look for interfaces

Flags:
`, cmd, cmd, cmd)
	for _, fs := range fss {
		fs.PrintDefaults()
	}
}
