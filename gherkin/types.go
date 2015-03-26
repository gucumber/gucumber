package gherkin

import "reflect"

// Feature represents the top-most construct in a Gherkin document. A feature
// contains one or more scenarios, which in turn contains multiple steps.
type Feature struct {
	// The filename where the feature was defined
	Filename string

	// The line number where the feature was defined
	Line int

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

	// The longest line length in the feature (including title)
	longestLine int
}

// Scenario represents a scenario (or background) of a given feature.
type Scenario struct {
	// The filename where the scenario was defined
	Filename string

	// The line number where the scenario was defined
	Line int

	// The scenario's title. For backgrounds, this is the empty string.
	Title string

	// Any tags associated with this scenario.
	Tags []Tag

	// All steps associated with the scenario.
	Steps []Step

	// Contains all scenario outline example data, if provided.
	Examples TabularDataMap

	// The longest line length in the scenario (including title)
	longestLine int
}

// Step represents an individual step making up a gucumber scenario.
type Step struct {
	// The filename where the step was defined
	Filename string

	// The line number where the step was defined
	Line int

	// The step's "type" (Given, When, Then, And, ...)
	//
	// Note that this field is normalized to the English form (e.g., "Given").
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

// TabularDataMap is tabular text data attached to a step organized in map
// form of the header name and its associated row data.
type TabularDataMap map[string][]string

// StepType represents a given step type.
type StepType string

// Tag is a string representation of a tag used in Gherkin syntax.
type Tag string

// ToMap converts a regular table to a map of header names to their row data.
// For example:
//
//     t := TabularData{[]string{"header1", "header2"}, []string{"col1", "col2"}}
//     t.ToMap()
//     // Output:
//     //   map[string][]string{
//     //     "header1": []string{"col1"},
//     //     "header2": []string{"col2"},
//     //   }
func (t TabularData) ToMap() TabularDataMap {
	m := TabularDataMap{}
	if len(t) > 1 {
		for _, th := range t[0] {
			m[th] = []string{}
		}
		for _, tr := range t[1:] {
			for c, td := range tr {
				m[t[0][c]] = append(m[t[0][c]], td)
			}
		}
	}
	return m
}

// NumRows returns the number of rows in a table map
func (t TabularDataMap) NumRows() int {
	if len(t) == 0 {
		return 0
	}
	return len(t[reflect.ValueOf(t).MapKeys()[0].String()])
}

func (s *Scenario) LongestLine() int {
	if s.longestLine == 0 {
		s.longestLine = len("Scenario: " + s.Title)
		for _, step := range s.Steps {
			if l := len(string(step.Type) + " " + step.Text); l > s.longestLine {
				s.longestLine = l
			}
		}
	}
	return s.longestLine
}

func (f *Feature) LongestLine() int {
	if f.longestLine == 0 {
		f.longestLine = len("Feature: " + f.Title)
		for _, s := range f.Scenarios {
			if l := s.LongestLine(); l > f.longestLine {
				f.longestLine = l
			}
		}
	}
	return f.longestLine
}
