package main

import (
	"io"
	"os"
)

type Config struct {
	Root    string
	DryRun  bool
	Verbose bool
	Stdout  io.Writer
}

func readConfig() Config {
	config := Config{
		Root:    "fs",
		DryRun:  false,
		Verbose: true,
		Stdout:  os.Stdout,
	}
	return config
}
