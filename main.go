package main

import (
	"github.com/gizo-network/gizo/cli"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	cli.Execute()
}
