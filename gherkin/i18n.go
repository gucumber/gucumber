package gherkin

type Language string

type Translation struct {
	Feature    string
	Background string
	Scenario   string
	And        string
	Given      string
	When       string
	Then       string
	Examples   string
}

const (
	LANG_EN = Language("en")
)

var (
	Translations = map[Language]Translation{
		LANG_EN: Translation{
			Feature:    "Feature",
			Background: "Background",
			Scenario:   "Scenario",
			And:        "And",
			Given:      "Given",
			When:       "When",
			Then:       "Then",
			Examples:   "Examples",
		},
	}
)
