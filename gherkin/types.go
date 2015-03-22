package gherkin

type Feature struct {
	Title       string
	Description string
	Tags        []Tag
	Background  Scenario
	Scenarios   []Scenario
}

type Scenario struct {
	Title    string
	Tags     []Tag
	Steps    []Step
	Examples interface{}
}

type Step struct {
	Type     StepType
	Text     string
	Argument interface{}
}

type StringData string

type TabularData [][]string

type StepType string

type Tag string
