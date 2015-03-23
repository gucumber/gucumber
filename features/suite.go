package main

import (
	. "github.com/lsegal/go-cucumber"
	"github.com/stretchr/testify/assert"
)

func main() {
	executions := 0

	Given(`^I have an initial step$`, func() {
		assert.Equal(T, 1, 1)
	})

	And(`^I have a second step$`, func() {
		assert.Equal(T, 2, 2)
	})

	When(`^I run the "(.+?)" command$`, func(s1 string) {
		assert.Equal(T, "cucumber.go", s1)
	})

	Then(`^this scenario should execute (\d+) time and pass$`, func(i1 int) {
		executions++
		assert.Equal(T, executions, i1)
	})

	RunMain()
}
