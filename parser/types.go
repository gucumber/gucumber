package parser

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
	Examples StepArgument
}

type Step struct {
	Type     StepType
	Text     string
	Argument StepArgument
}

type StepType string

type StepArgument string

type Tag string
