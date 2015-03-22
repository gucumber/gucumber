package cucumber_test

import (
	"testing"

	. "github.com/lsegal/go-cucumber"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSteps(t *testing.T) {
	count := 0
	str := ""
	fl := 0.0

	Given(`^I have a test with (\d+)$`, func(i int) { count += i })
	When(`^I have a condition of (\d+) with decimal (-?\d+\.\d+)$`, func(i int64, f float64) { count += int(i); fl = f })
	And(`^I have another condition with "(.+?)"$`, func(s string) { str = s })
	Then(`^something will happen with text$`, func(data string) { str += data })

	RegisteredSteps.Execute("I have a test with 3", "")
	RegisteredSteps.Execute("I have a condition of 5 with decimal -3.14159", "")
	RegisteredSteps.Execute("I have another condition with \"arbitrary text\"", "")
	RegisteredSteps.Execute("something will happen with text", " and hello world")

	assert.Equal(t, 8, count)
	assert.Equal(t, "arbitrary text and hello world", str)
	assert.Equal(t, -3.14159, fl)
}
