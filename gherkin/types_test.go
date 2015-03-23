package gherkin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTabularDataToMap(t *testing.T) {
	tab := TabularData{
		[]string{"a", "b", "c", "d"},
		[]string{"1", "2", "3", "4"},
		[]string{"5", "6", "7", "8"},
		[]string{"9", "A", "B", "C"},
	}

	m := TabularDataMap{
		"a": []string{"1", "5", "9"},
		"b": []string{"2", "6", "A"},
		"c": []string{"3", "7", "B"},
		"d": []string{"4", "8", "C"},
	}

	assert.Equal(t, m, tab.ToMap())
}

func TestTabularDataMapEmpty(t *testing.T) {
	var tab TabularData
	var m TabularDataMap

	// only headers
	tab = TabularData{[]string{"a", "b", "c", "d"}}
	m = TabularDataMap{}
	assert.Equal(t, m, tab.ToMap())

	// completely empty
	tab = TabularData{}
	m = TabularDataMap{}
	assert.Equal(t, m, tab.ToMap())
}
