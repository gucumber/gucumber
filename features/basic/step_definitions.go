package basic

import (
	. "github.com/lsegal/gucumber"
	"github.com/stretchr/testify/assert"
)

func init() {
	executions := 0

	Given(`^I have an initial step$`, func() {
		assert.Equal(T, 1, 1)
	})

	And(`^I have a second step$`, func() {
		assert.Equal(T, 2, 2)
	})

	When(`^I run the "(.+?)" command$`, func(s1 string) {
		assert.Equal(T, "gucumber.go", s1)
	})

	Then(`^this scenario should execute (\d+) time and pass$`, func(i1 int) {
		executions++
		assert.Equal(T, executions, i1)
	})
}
