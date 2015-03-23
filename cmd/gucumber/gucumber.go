package main

import (
	"os"

	"github.com/lsegal/gucumber"
)

func main() {
	gucumber.RunMain(os.Args[1:]...)
}
