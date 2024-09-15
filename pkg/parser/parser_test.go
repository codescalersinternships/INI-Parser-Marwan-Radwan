package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.data

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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.data

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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.data

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

	t.Run("Invalid Key-Value Without Separator", func(t *testing.T) {

		input := `
[section1]
key1=value1
key2
key3=value3
`

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()

		err := p.parse(scanner)

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

func TestGetSectionNames(t *testing.T) {
	t.Run("Single Section", func(t *testing.T) {
		input := `
[section1]
key1=value1
key2=value2
`
		expected := []string{"section1"}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()

		err := p.parse(scanner)
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		sectionNames := p.GetSectionNames()

		if len(sectionNames) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(sectionNames))
		}

		if !reflect.DeepEqual(sectionNames, expected) {
			t.Errorf("Expected section names %v, got %v", expected, sectionNames)
		}

	})

	t.Run("Multiple Section", func(t *testing.T) {
		input := `
[section1]
key1=value1
key2=value2

[section2]
keyA=valueA
keyB=valueB
`
		expected := []string{"section1", "section2"}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		sectionNames := p.GetSectionNames()

		if len(sectionNames) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(sectionNames))
		}

		if !reflect.DeepEqual(sectionNames, expected) {
			t.Errorf("Expected section names %v, got %v", expected, sectionNames)
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
		expected := []string{"section1", "section2"}

		scanner := bufio.NewScanner(strings.NewReader(input))

		p := NewParser()
		err := p.parse(scanner)

		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		sectionNames := p.GetSectionNames()

		if len(sectionNames) != len(expected) {
			t.Errorf("Expected %d sections, got %d", len(expected), len(sectionNames))
		}

		if !reflect.DeepEqual(sectionNames, expected) {
			t.Errorf("Expected section names %v, got %v", expected, sectionNames)
		}
	})

}

func TestGetSections(t *testing.T) {
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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.GetSections()

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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.GetSections()

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
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.GetSections()

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

func TestGet(t *testing.T) {
	t.Run("Single Section", func(t *testing.T) {
		input := `
[section1]
key1=value1
key2=value2
`
		expected := "value1"

		p := NewParser()

		scanner := bufio.NewScanner(strings.NewReader(input))

		err := p.parse(scanner)
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		val, ok := p.Get("section1", "key1")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}

		val2, ok := p.Get("section1", "key3")
		if !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", "", false, val2, ok)
		}

	})

	t.Run("Multiple Section", func(t *testing.T) {
		input := `
[section1]
key1=value1
key2=value2
[section2]
key3=value3
key4=value4
`
		expected := "value1"

		p := NewParser()

		scanner := bufio.NewScanner(strings.NewReader(input))

		err := p.parse(scanner)
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		val, ok := p.Get("section1", "key1")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}

		val2, ok := p.Get("section1", "key3")
		if !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", "", false, val2, ok)
		}

	})
}

func TestSet(t *testing.T) {
	t.Run("Set Key in Existing Section", func(t *testing.T) {
		p := NewParser()
		p.Set("section1", "key1", "value1")

		expected := "value1"
		val, ok := p.Get("section1", "key1")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}
	})

	t.Run("Set Key in New Section", func(t *testing.T) {
		p := NewParser()
		p.Set("section2", "key2", "value2")

		expected := "value2"
		val, ok := p.Get("section2", "key2")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}
	})

	t.Run("Overwrite Existing Key", func(t *testing.T) {
		p := NewParser()
		p.Set("section1", "key1", "value1")
		p.Set("section1", "key1", "newValue1")

		expected := "newValue1"
		val, ok := p.Get("section1", "key1")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}
	})

	t.Run("Set Key in Empty Section", func(t *testing.T) {
		p := NewParser()
		p.Set("section3", "key3", "value3")

		expected := "value3"
		val, ok := p.Get("section3", "key3")
		if val != expected || !ok {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected, true, val, ok)
		}
	})

	t.Run("Set Multiple Keys in Same Section", func(t *testing.T) {
		p := NewParser()
		p.Set("section4", "key4", "value4")
		p.Set("section4", "key5", "value5")

		expected1 := "value4"
		expected2 := "value5"

		val1, ok1 := p.Get("section4", "key4")
		if val1 != expected1 || !ok1 {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected1, true, val1, ok1)
		}

		val2, ok2 := p.Get("section4", "key5")
		if val2 != expected2 || !ok2 {
			t.Errorf("Expected to get %s and %v, got %s and %v", expected2, true, val2, ok2)
		}
	})
}

func TestToString(t *testing.T) {
	t.Run("Single Section", func(t *testing.T) {
		p := NewParser()

		p.data = map[string]map[string]string{
			"section1": {
				"key1": "value1",
				"key2": "value2",
			},
		}
		p.sections = []string{"section1"}

		expected := `[section1]
key1=value1
key2=value2
`

		result := p.ToString()
		if result != expected {
			t.Errorf("Expected to get \n%s\n ,but got:\n%s", expected, result)
		}
	})

	t.Run("Multiple Sections", func(t *testing.T) {
		p := NewParser()

		p.data = map[string]map[string]string{
			"section1": {
				"key1": "value1",
				"key2": "value2",
			},
			"section2": {
				"keyA": "valueA",
				"keyB": "valueB",
			},
		}
		p.sections = []string{"section1", "section2"}

		expected := `[section1]
key1=value1
key2=value2
[section2]
keyA=valueA
keyB=valueB
`

		result := p.ToString()
		if result != expected {
			t.Errorf("Expected to get \n%s\n ,but got:\n%s", expected, result)
		}
	})

	t.Run("No Sections", func(t *testing.T) {
		p := NewParser()

		expected := ``

		result := p.ToString()
		if result != expected {
			t.Errorf("Expected to get \n%s\n ,but got:\n%s", expected, result)
		}
	})
}

func TestLoadFromString(t *testing.T) {
	t.Run("Valid INI String", func(t *testing.T) {
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

		p := NewParser()
		err := p.LoadFromString(input)

		if err != nil {
			t.Errorf("Error loading from string: %v", err)
		}

		iniData := p.data

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

	t.Run("Empty String", func(t *testing.T) {
		input := ``

		p := NewParser()
		err := p.LoadFromString(input)

		if err != nil {
			t.Errorf("Error loading from string: %v", err)
		}

		iniData := p.data

		if len(iniData) != 0 {
			t.Errorf("Expected 0 sections, got %d", len(iniData))
		}
	})

	t.Run("String with Comments and Empty Lines", func(t *testing.T) {
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

		p := NewParser()
		err := p.LoadFromString(input)

		if err != nil {
			t.Errorf("Error loading from string: %v", err)
		}

		iniData := p.data

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

	t.Run("Invalid INI String", func(t *testing.T) {
		input := `
[section1]
key1=value1
key2
key3=value3
`
		p := NewParser()
		err := p.LoadFromString(input)

		if err == nil {
			t.Errorf("Expected error for invalid INI string, but got no error")
		} else {
			expectedError := "line 4: invalid key-value pair: key2"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})
}

func TestParseFile(t *testing.T) {
	t.Run("Valid INI File", func(t *testing.T) {
		content := `
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

		// Create a temporary directory and .ini file
		dir := t.TempDir()
		filePath := filepath.Join(dir, "test.ini")
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create temporary INI file: %v", err)
		}

		p := NewParser()

		err = p.ParseFile(filePath)
		if err != nil {
			t.Errorf("Error parsing ini file: %v", err)
		}

		iniData := p.data

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

	t.Run("Invalid File Extension", func(t *testing.T) {

		content := `
		[section1]
		key1=value1
		key2=value2
		`

		// Create a temporary file with an invalid extension
		dir := t.TempDir()
		filePath := filepath.Join(dir, "test.txt")
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}

		p := NewParser()

		err = p.ParseFile(filePath)
		if err == nil {
			t.Errorf("Expected error for invalid file extension, but got no error")
		} else {
			expectedError := ".ini format is only support"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		p := NewParser()

		err := p.ParseFile("non_existent.ini")
		if err == nil {
			t.Errorf("Expected error for non-existent file, but got no error")
		} else {
			expectedError := "open non_existent.ini: no such file or directory"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
			}
		}
	})
}
