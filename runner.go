package cucumber

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lsegal/go-cucumber/gherkin"
	"github.com/shiena/ansicolor"
)

const (
	clrRed    = "31"
	clrGreen  = "32"
	clrYellow = "33"
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

func (c *Context) RunDir(dir string) error {
	g, err := filepath.Glob(filepath.Join(dir, "*.feature"))
	if err != nil {
		return err
	}

	g2, err := filepath.Glob(filepath.Join(dir, "**", "*.feature"))
	if err != nil {
		return err
	}

	g = append(g, g2...)
	return c.RunFiles(g)
}

func (c *Context) RunFiles(files []string) error {
	r := Runner{
		Context:   c,
		Features:  []*gherkin.Feature{},
		Results:   []RunnerResult{},
		Unmatched: []*gherkin.Step{},
	}

	for _, file := range files {
		fd, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fd.Close()

		b, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}

		fs, err := gherkin.ParseFilename(string(b), file)
		if err != nil {
			return err
		}

		for _, f := range fs {
			r.Features = append(r.Features, &f)
		}
	}

	r.run()
	return nil
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
						skipping = true
						clr = clrRed
					}
					done <- true
				}()

				b, err := c.Execute(t, step.Text, step.Argument)
				if err != nil {
					t.Error(err)
				}
				if !b {
					c.Unmatched = append(c.Unmatched, &step)
					t.Skip("no match function for step")
				}
			}()
			<-done
		}

		c.line(clr, "    %s %s", step.Type, step.Text)
	}
}

func (c *Runner) line(clr, text string, args ...interface{}) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	fmt.Fprintf(w, "\033[%sm%s\033[0;0m\n", clr, fmt.Sprintf(text, args...))
}
