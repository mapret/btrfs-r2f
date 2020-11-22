package main

type Config struct {
	root   string
	dryRun bool
}

func readConfig() Config {
	config := Config{
		root:   "fs",
		dryRun: false,
	}
	return config
}
