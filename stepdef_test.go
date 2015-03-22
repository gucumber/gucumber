package cucumber_test

import (
	"testing"

	. "github.com/lsegal/go-cucumber"
	"github.com/lsegal/go-cucumber/gherkin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSteps(t *testing.T) {
	count := 0
	str := ""
	fl := 0.0

	Given(`^I have a test with (\d+)$`, func(i int) { count += i })
	When(`^I have a condition of (\d+) with decimal (-?\d+\.\d+)$`, func(i int64, f float64) { count += int(i); fl = f })
	And(`^I have another condition with "(.+?)"$`, func(s string) { str = s })
	Then(`^something will happen with text$`, func(data string) { str += " " + data })
	And(`^something will happen with a table:$`, func(table gherkin.TabularData) {
		str += " " + table[0][0] + table[0][1] + table[1][0] + table[1][1]
	})
	And(`^something will happen with a table:$`, func(table [][]string) {
		str += " " + table[0][0] + table[0][1] + table[1][0] + table[1][1]
	})
	And(`^I can pass in test (.+?)$`, func(tt *testing.T, data string) {
		assert.Equal(t, t, tt)
		str += " " + data
	})

	GlobalContext.Execute(t, "I have a test with 3", "")
	GlobalContext.Execute(t, "I have a condition of 5 with decimal -3.14159", "")
	GlobalContext.Execute(t, "I have another condition with \"arbitrary text\"", "")
	GlobalContext.Execute(t, "something will happen with text", "and hello world")
	GlobalContext.Execute(t, "something will happen with a table:",
		gherkin.TabularData{[]string{"a", "b"}, []string{"c", "d"}})
	GlobalContext.Execute(t, "I can pass in test context", "")

	assert.Equal(t, 8, count)
	assert.Equal(t, "arbitrary text and hello world abcd abcd context", str)
	assert.Equal(t, -3.14159, fl)
}
