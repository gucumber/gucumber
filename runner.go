package gucumber

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/lsegal/gucumber/gherkin"
	"github.com/shiena/ansicolor"
)

const (
	clrRed    = "31"
	clrGreen  = "32"
	clrYellow = "33"

	txtUnmatchInt   = `(\d+)`
	txtUnmatchFloat = `(-?\d+(?:\.\d+)?)`
	txtUnmatchStr   = `"(.+?)"`
)

var (
	reUnmatchInt   = regexp.MustCompile(txtUnmatchInt)
	reUnmatchFloat = regexp.MustCompile(txtUnmatchFloat)
	reUnmatchStr   = regexp.MustCompile(txtUnmatchStr)
)

type Runner struct {
	*Context
	Features  []*gherkin.Feature
	Results   []RunnerResult
	Unmatched []*gherkin.Step
	FailCount int
	SkipCount int
}

type RunnerResult struct {
	*testing.T
	*gherkin.Feature
	*gherkin.Scenario
}

func (c *Context) RunDir(dir string) (*Runner, error) {
	g, err := filepath.Glob(filepath.Join(dir, "*.feature"))
	if err != nil {
		return nil, err
	}

	g2, err := filepath.Glob(filepath.Join(dir, "**", "*.feature"))
	if err != nil {
		return nil, err
	}

	g = append(g, g2...)
	return c.RunFiles(g)
}

func (c *Context) RunFiles(files []string) (*Runner, error) {
	r := Runner{
		Context:   c,
		Features:  []*gherkin.Feature{},
		Results:   []RunnerResult{},
		Unmatched: []*gherkin.Step{},
	}

	for _, file := range files {
		fd, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fd.Close()

		b, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}

		fs, err := gherkin.ParseFilename(string(b), file)
		if err != nil {
			return nil, err
		}

		for _, f := range fs {
			r.Features = append(r.Features, &f)
		}
	}

	r.run()
	return &r, nil
}

func (c *Runner) MissingMatcherStubs() string {
	var buf bytes.Buffer

	buf.WriteString(`import . "github.com/lsegal/gucumber"` + "\n\n")
	buf.WriteString("func init() {\n")

	for _, m := range c.Unmatched {
		numInts, numFloats, numStrs := 1, 1, 1
		str, args := m.Text, []string{}
		str = reUnmatchInt.ReplaceAllStringFunc(str, func(s string) string {
			args = append(args, fmt.Sprintf("i%d int", numInts))
			numInts++
			return txtUnmatchInt
		})
		str = reUnmatchFloat.ReplaceAllStringFunc(str, func(s string) string {
			args = append(args, fmt.Sprintf("s%d float64", numFloats))
			numFloats++
			return txtUnmatchFloat
		})
		str = reUnmatchStr.ReplaceAllStringFunc(str, func(s string) string {
			args = append(args, fmt.Sprintf("s%d string", numStrs))
			numStrs++
			return txtUnmatchStr
		})

		switch m.Argument.(type) {
		case gherkin.TabularData:
			args = append(args, "table [][]string")
		case gherkin.StringData:
			args = append(args, "data string")
		}

		fmt.Fprintf(&buf, "\t%s(`^%s$`, func(%s) {\n\t\tT.Skip() // pending\n\t})\n\n",
			m.Type, str, strings.Join(args, ", "))
	}

	buf.WriteString("}\n")
	return buf.String()
}

func (c *Runner) run() {
	for _, f := range c.Features {
		c.runFeature(f)
	}

	c.line("0;1", "finished (%d passed, %d failed, %d skipped).\n",
		len(c.Results)-c.FailCount-c.SkipCount, c.FailCount, c.SkipCount)
}

func (c *Runner) runFeature(f *gherkin.Feature) {
	c.line("0;1", "Feature: %s", f.Title)

	if f.Background.Steps != nil {
		c.line("0;1", "  Background:")
		c.runScenario(f, &f.Background)
		c.line("0", "")
	}

	for _, s := range f.Scenarios {
		c.line("0;1", "  Scenario: %s", s.Title)
		c.runScenario(f, &s)
		c.line("0", "")
	}
}

func (c *Runner) runScenario(f *gherkin.Feature, s *gherkin.Scenario) {
	t := &testing.T{}
	skipping := false
	clr := clrGreen

	for _, step := range s.Steps {
		found := false
		if !skipping {
			done := make(chan bool)
			go func() {
				defer func() {
					c.Results = append(c.Results, RunnerResult{t, f, s})

					if t.Skipped() {
						c.SkipCount++
						skipping = true
						clr = clrYellow
					} else if t.Failed() {
						c.FailCount++
						clr = clrRed
					}
					done <- true
				}()

				f, err := c.Execute(t, step.Text, step.Argument)
				if err != nil {
					t.Error(err)
				}
				found = f

				if !f {
					t.Skip("no match function for step")
				}
			}()
			<-done
		}

		if skipping && !found {
			cstep := step
			c.Unmatched = append(c.Unmatched, &cstep)
		}

		c.line(clr, "    %s %s", step.Type, step.Text)
	}
}

func (c *Runner) line(clr, text string, args ...interface{}) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	fmt.Fprintf(w, "\033[%sm%s\033[0;0m\n", clr, fmt.Sprintf(text, args...))
}
