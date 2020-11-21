package main

type Config struct {
	root string
}

func readConfig() Config {
	config := Config{root: "fs"}
	return config
}
