package parser

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("Single Section", func(t *testing.T) {

		input := `
[section1]
key1=value1
key2=value2
`
		expected := map[string]map[string]string{
			"section1": {
				"key1": "value1",
				"key2": "value2",
			},
		}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		iniData := p.data

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		if len(iniData) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(iniData))
		}

		for section, keys := range expected {
			if _, ok := iniData[section]; !ok {
				t.Errorf("Expected section %s, not found", section)
			}

			for key, expectedValue := range keys {
				if val, ok := iniData[section][key]; !ok || val != expectedValue {
					t.Errorf("Expected key %s in section %s to have value %s, got %s", key, section, expectedValue, val)
				}
			}
		}
	})

	t.Run("Multiple Sections", func(t *testing.T) {

		input := `
[section1]
key1=value1
key2=value2

[section2]
keyA=valueA
keyB=valueB
`
		expected := map[string]map[string]string{
			"section1": {
				"key1": "value1",
				"key2": "value2",
			},
			"section2": {
				"keyA": "valueA",
				"keyB": "valueB",
			},
		}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		iniData := p.data

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		if len(iniData) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(iniData))
		}

		for section, keys := range expected {
			if _, ok := iniData[section]; !ok {
				t.Errorf("Expected section %s, not found", section)
			}

			for key, expectedValue := range keys {
				if val, ok := iniData[section][key]; !ok || val != expectedValue {
					t.Errorf("Expected key %s in section %s to have value %s, got %s", key, section, expectedValue, val)
				}
			}
		}
	})

	t.Run("Comments and Empty Lines", func(t *testing.T) {

		input := `
; this is a comment
# this is another comment

[section1]
key1=value1

[section2]
keyA=valueA		
`
		expected := map[string]map[string]string{
			"section1": {
				"key1": "value1",
			},
			"section2": {
				"keyA": "valueA",
			},
		}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		iniData := p.data

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		if len(iniData) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(iniData))
		}

		for section, keys := range expected {
			if _, ok := iniData[section]; !ok {
				t.Errorf("Expected section %s, not found", section)
			}

			for key, expectedValue := range keys {
				if val, ok := iniData[section][key]; !ok || val != expectedValue {
					t.Errorf("Expected key %s in section %s to have value %s, got %s", key, section, expectedValue, val)
				}
			}
		}
	})

	t.Run("Global Keys - No Sections", func(t *testing.T) {

		input := `
key1=value1
key2=value2		
`
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		iniData := p.globalKeys

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		if len(iniData) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(iniData))
		}

		for key, expectedValue := range iniData {
			if val, ok := iniData[key]; !ok || val != expectedValue {
				t.Errorf("Expected key %s to have value %s, got %s", key, expectedValue, val)
			}
		}
	})

	t.Run("Global Keys - With Sections", func(t *testing.T) {

		input := `
key1=value1
key2=value2

[section1]
keyA=valueA
`
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		iniData := p.globalKeys

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		if len(iniData) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(iniData))
		}

		for key, expectedValue := range iniData {
			if val, ok := iniData[key]; !ok || val != expectedValue {
				t.Errorf("Expected key %s to have value %s, got %s", key, expectedValue, val)
			}
		}
	})

	t.Run("Invalid Key-Value Without Separator", func(t *testing.T) {

		input := `
[section1]
key1=value1
key2
key3=value3
`

		// Create a scanner to read the input
		scanner := bufio.NewScanner(strings.NewReader(input))

		// Initialize the parser
		p := NewParser()

		// Parse the input and expect an error
		err := p.parse(scanner)

		// We expect an error, so check that it is not nil
		if err == nil {
			t.Errorf("Expected error for invalid key-value pairs, but got no error")
		} else {
			expectedError := "line 4: invalid key-value pair: key2"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})

	t.Run("Invalid Key-Value With Separator", func(t *testing.T) {
		input := `
[section1]
key1=value1
keyWithSeparatorButNoValue =
= valueWithoutKey
key3=value3
	`
		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()

		err := p.parse(scanner)

		if err == nil {
			t.Errorf("Expected error for invalid key-value pairs with separator, but got no error")
		} else {
			expectedError := "line 4: value cannot be empty: keyWithSeparatorButNoValue ="
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})

	t.Run("Malformed Sections", func(t *testing.T) {

		input := `
[section1
key1=value1
key2=value2

section2]
keyA=valueA
keyB=valueB

section3]
key3 = value3

[anotherInvalidSection
key4 = value4

[]
key5 = value5
`
		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		if err == nil {
			t.Errorf("Expected error for invalid key-value pairs with separator, but got no error")
		} else {
			expectedError := "line 2: invalid key-value pair: [section1"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})

	t.Run("Backslash", func(t *testing.T) {
		input := `
key=long_value \
that_spans_multiple_lines
`

		fmt.Println(input)
	})
}
