package main

import (
	"os"

	"github.com/lsegal/go-cucumber"
)

func main() {
	err := cucumber.GlobalContext.RunDir(os.Args[1])
	if err != nil {
		panic(err)
	}
}
