package cucumber

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

var (
	GlobalContext = Context{Steps: []StepDefinition{}}
	T             *testing.T
)

func Given(match string, fn interface{}) {
	GlobalContext.Given(match, fn)
}

func Then(match string, fn interface{}) {
	GlobalContext.Then(match, fn)
}

func When(match string, fn interface{}) {
	GlobalContext.When(match, fn)
}

func And(match string, fn interface{}) {
	GlobalContext.And(match, fn)
}

type Context struct {
	Steps []StepDefinition
	T     *testing.T
}

func (c *Context) addStep(match string, fn interface{}) {
	c.Steps = append(c.Steps, StepDefinition{
		Matcher:  regexp.MustCompile(match),
		Function: reflect.ValueOf(fn),
	})
}

func (c *Context) Given(match string, fn interface{}) {
	c.addStep(match, fn)
}

func (c *Context) Then(match string, fn interface{}) {
	c.addStep(match, fn)
}

func (c *Context) When(match string, fn interface{}) {
	c.addStep(match, fn)
}

func (c *Context) And(match string, fn interface{}) {
	c.addStep(match, fn)
}

func (c *Context) Execute(line string, arg interface{}) error {
	for _, step := range c.Steps {
		step.CallIfMatch(line, arg)
	}
	return nil
}

func (c *Context) ExecuteTest(t *testing.T, line string, arg interface{}) {
	T = t
	c.Execute(line, arg)
}

type StepDefinition struct {
	Matcher  *regexp.Regexp
	Function reflect.Value
}

func (s *StepDefinition) CallIfMatch(line string, arg interface{}) {
	if match := s.Matcher.FindStringSubmatch(line); match != nil {
		match = match[1:] // discard full line match

		// adjust arity if there is step arg data
		numArgs := len(match)
		if arg != "" {
			numArgs++
		}

		t := s.Function.Type()
		if numArgs != t.NumIn() { // function has different arity
			return // TODO raise error
		}

		values := make([]reflect.Value, numArgs)
		for i := 0; i < t.NumIn(); i++ {
			var v interface{}
			switch t.In(i).Kind() {
			case reflect.Int:
				i, _ := strconv.ParseInt(match[i], 10, 32)
				v = int(i)
			case reflect.Int64:
				v, _ = strconv.ParseInt(match[i], 10, 64)
			case reflect.String:
				// this could be from `arg`, check match index
				if i >= len(match) {
					v = arg
				} else {
					v = match[i]
				}
			case reflect.Float64:
				v, _ = strconv.ParseFloat(match[i], 64)
			default:
				panic("type " + t.String() + "is not supported.")
			}

			values[i] = reflect.ValueOf(v)
		}

		s.Function.Call(values)
	}
}
