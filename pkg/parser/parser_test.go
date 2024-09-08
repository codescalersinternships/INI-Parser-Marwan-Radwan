package parser

import (
	"bufio"
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
		iniData, err := p.parse(scanner)

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
		iniData, err := p.parse(scanner)

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
		iniData, err := p.parse(scanner)

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

	t.Run("No Sections", func(t *testing.T) {

		input := `
key1=value1
key2=value2		
`
		expected := map[string]map[string]string{}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		iniData, err := p.parse(scanner)

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

}
