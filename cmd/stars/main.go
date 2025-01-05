package main

import (
	"log"

	"github.com/jiangjiax/stars/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
