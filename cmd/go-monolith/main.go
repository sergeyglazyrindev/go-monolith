package main

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := gomonolith.NewApp(environment)
	app1.ExecuteCommand()
}
