package gherkin

// Feature represents the top-most construct in a Gherkin document. A feature
// contains one or more scenarios, which in turn contains multiple steps.
type Feature struct {
	// The feature's title.
	Title string

	// A longer description of the feature. This is not used during runtime.
	Description string

	// Any tags associated with this feature.
	Tags []Tag

	// Any background scenario data that is executed prior to scenarios.
	Background Scenario

	// The scenarios associated with this feature.
	Scenarios []Scenario
}

// Scenario represents a scenario (or background) of a given feature.
type Scenario struct {
	// The scenario's title. For backgrounds, this is the empty string.
	Title string

	// Any tags associated with this scenario.
	Tags []Tag

	// All steps associated with the scenario.
	Steps []Step

	// Contains all scenario outline example data, if provided.
	Examples TabularData
}

// Step represents an individual step making up a Cucumber scenario.
type Step struct {
	// The step's "type" (Given, When, Then, And, ...)
	//
	// Note that this field is not normalized across languages and represents
	// the type as it appeared in the original source text. If you want to
	// check against a set of known types ("Given" specifically, for example),
	// you should compare against names found in Translations["LANG"] instead.
	Type StepType

	// The text contained in the step (minus the "Type" prefix).
	Text string

	// Argument represents multi-line argument data attached to a step.
	// This value is an interface{} but is only ever set to a StringData
	// or TabularData type.
	Argument interface{}
}

// StringData is multi-line docstring text attached to a step.
type StringData string

// TabularData is tabular text data attached to a step.
type TabularData [][]string

// StepType represents a given step type.
type StepType string

// Tag is a string representation of a tag used in Gherkin syntax.
type Tag string
