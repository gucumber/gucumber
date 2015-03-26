package gucumber

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

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
	reUnmatchStr   = regexp.MustCompile(`(<|").+?("|>)`)
	reOutlineVal   = regexp.MustCompile(`<(.+?)>`)
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
	*TestingT
	*gherkin.Feature
	*gherkin.Scenario
}

func (c *Context) RunDir(dir string) (*Runner, error) {
	g, _ := filepath.Glob(filepath.Join(dir, "*.feature"))
	g2, _ := filepath.Glob(filepath.Join(dir, "**", "*.feature"))
	g = append(g, g2...)

	runner, err := c.RunFiles(g)
	if err != nil {
		panic(err)
	}

	if len(runner.Unmatched) > 0 {
		fmt.Println("Some steps were missing, you can add them by using the following step definition stubs: ")
		fmt.Println("")
		fmt.Print(runner.MissingMatcherStubs())
	}

	return runner, err
}

func (c *Context) RunFiles(featureFiles []string) (*Runner, error) {
	r := Runner{
		Context:   c,
		Features:  []*gherkin.Feature{},
		Results:   []RunnerResult{},
		Unmatched: []*gherkin.Step{},
	}

	for _, file := range featureFiles {
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
	matches := map[string]bool{}

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

		// Don't duplicate matchers. This is mostly for scenario outlines.
		if matches[str] {
			continue
		}
		matches[str] = true

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

	os.Exit(c.FailCount)
}

func (c *Runner) runFeature(f *gherkin.Feature) {
	c.line("0;1", "Feature: %s", f.Title)

	if f.Background.Steps != nil {
		c.runScenario("Background", f, &f.Background)
	}

	for _, s := range f.Scenarios {
		c.runScenario("Scenario", f, &s)
	}
}

func (c *Runner) runScenario(title string, f *gherkin.Feature, s *gherkin.Scenario) {
	if len(s.Examples) > 1 { // run scenario outline data
		for i, rows := 0, s.Examples.NumRows(); i < rows; i++ {
			other := gherkin.Scenario{
				Filename: s.Filename,
				Line:     s.Line,
				Title:    s.Title,
				Examples: gherkin.TabularDataMap{},
				Steps:    []gherkin.Step{},
			}

			for _, step := range s.Steps {
				step.Text = reOutlineVal.ReplaceAllStringFunc(step.Text, func(t string) string {
					return s.Examples[t[1:len(t)-1]][i]
				})
				other.Steps = append(other.Steps, step)
			}
			c.runScenario(title, f, &other)
		}
		return
	}

	t := &TestingT{}
	skipping := false
	clr := clrGreen

	c.fileLine("0;1", "  %s: %s", s.Filename, s.Line, s.LongestLine(),
		title, s.Title)

	for _, step := range s.Steps {
		errCount := len(t.errors)
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

		c.fileLine(clr, "    %s %s", step.Filename, step.Line,
			s.LongestLine(), step.Type, step.Text)

		if len(t.errors) > errCount {
			c.line(clrRed, "\n"+t.errors[len(t.errors)-1].message)
		}
	}
	c.line("0", "")
}

var writer = ansicolor.NewAnsiColorWriter(os.Stdout)

func (c *Runner) line(clr, text string, args ...interface{}) {
	fmt.Fprintf(writer, "\033[%sm%s\033[0;0m\n", clr, fmt.Sprintf(text, args...))
}

func (c *Runner) fileLine(clr, text, filename string, line int, max int, args ...interface{}) {
	space, str := "", fmt.Sprintf(text, args...)
	if l := max + 5 - len(str); l > 0 {
		space = strings.Repeat(" ", l)
	}
	comment := fmt.Sprintf("%s \033[39;0m# %s:%d", space, filename, line)
	c.line(clr, "%s%s", str, comment)
}

type Tester interface {
	Errorf(format string, args ...interface{})
}

type TestingT struct {
	skipped bool
	errors  []TestError
}

type TestError struct {
	message string
	stack   []byte
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	var buf bytes.Buffer

	str := fmt.Sprintf(format, args...)
	sbuf := make([]byte, 8192)
	for {
		size := runtime.Stack(sbuf, false)
		if size < len(sbuf) {
			break
		}
		buf.Write(sbuf[0:size])
	}

	t.errors = append(t.errors, TestError{message: str, stack: buf.Bytes()})
}

func (t *TestingT) Skip(args ...interface{}) {
	t.skipped = true
}

func (t *TestingT) Skipped() bool {
	return t.skipped
}

func (t *TestingT) Failed() bool {
	return len(t.errors) > 0
}

func (t *TestingT) Error(err error) {
	t.errors = append(t.errors, TestError{message: err.Error()})
}
