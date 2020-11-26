package main

import (
	"fmt"
	"io"
	"os"
)

type Config struct {
	Root      string
	DryRun    bool
	Verbose   bool
	Stdout    io.Writer
	InputPath string
}

var usageString = fmt.Sprintf(`Usage: %s [options...]

-i --input
    Command stream file, stdin is used if not supplied
-o --output
    Target directory, default value "btrfs-subvolume"
-d --dry-run
    Dont actually create files, only print commands (useful with -v)
-v --verbose
    Print every command that is processed to stdout
`, programName)

func printUsage() {
	fmt.Print(usageString)
}

func readConfig(args []string) Config {
	config := Config{
		Root:      "btrfs-subvolume",
		DryRun:    false,
		Verbose:   false,
		Stdout:    os.Stdout,
		InputPath: "",
	}

	for i := 0; i < len(args); i++ {
		if i == 0 {
		} else if args[i] == "-i" || args[i] == "--input" {
			i++
			config.InputPath = args[i]
		} else if args[i] == "-o" || args[i] == "output" {
			i++
			config.Root = args[i]
		} else if args[i] == "-d" || args[i] == "--dry-run" {
			config.DryRun = true
		} else if args[i] == "-v" || args[i] == "--verbose" {
			config.Verbose = true
		} else {
			fmt.Printf("Unknown argument: %s\n\n", args[i])
			printUsage()
			os.Exit(1)
		}
	}

	if !isStdinPipeConnected() {
		config.InputPath = args[len(args)-1]
	}

	return config
}
