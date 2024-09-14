package parser

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

type Parser struct {
	data       map[string]map[string]string
	globalKeys map[string]string
	sections   []string // Maintain the sections in order
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		data:       make(map[string]map[string]string),
		globalKeys: make(map[string]string),
		sections:   []string{},
	}
}

func (p *Parser) parse(scanner *bufio.Scanner) error {
	currentSection := ""
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			p.createSectionIfNotExist(currentSection)
		} else {
			err := p.parseKeyValue(line, lineNumber, currentSection)
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func (p *Parser) createSectionIfNotExist(section string) {
	if _, exist := p.data[section]; !exist {
		p.data[section] = make(map[string]string)
		p.sections = append(p.sections, section)
	}
}

func (p *Parser) parseKeyValue(line string, lineNumber int, currentSection string) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("line %d: invalid key-value pair: %s", lineNumber, line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return fmt.Errorf("line %d: key cannot be empty: %s", lineNumber, line)
	}
	if value == "" {
		return fmt.Errorf("line %d: value cannot be empty: %s", lineNumber, line)
	}

	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\r`, "\r")
	value = strings.ReplaceAll(value, `\t`, "\t")
	value = strings.Trim(value, `"`)

	if currentSection == "" {
		p.globalKeys[key] = value
	} else {
		p.data[currentSection][key] = value
	}

	return nil
}

// GetSectionNames lists all section names in the file.
func (p *Parser) GetSectionNames() []string {
	return p.sections
}

// GetSections returns a map of sections in the INI file, each section is represented by a map of key-value pairs.
func (p *Parser) GetSections() map[string]map[string]string {
	return p.data
}

// GetGlobalKeys returns a map of global keys in the parser.
func (p *Parser) GetGlobalKeys() map[string]string {
	return p.globalKeys
}

// Get retrieves the value associated with the given section and key from the files's data.
func (p *Parser) Get(section string, key string) (string, bool) {
	if sectionData, ok := p.data[section]; ok {
		return sectionData[key], ok
	}
	return "", false
}

// Set sets the value of a key in a specific section of the INI file.
func (p *Parser) Set(section string, key string, value string) {
	if _, exist := p.data[section]; !exist {
		p.data[section] = make(map[string]string)
	}
	p.data[section][key] = value
	p.sections = append(p.sections, section)
}

// ToString returns a string representation of the Parser object.
func (p *Parser) ToString() string {
	var builder strings.Builder

	// Write global keys first
	for k, v := range p.globalKeys {
		fmt.Fprintf(&builder, "%s=%s\n", k, v)
	}

	// Then write sections in order
	for _, sectionName := range p.sections {
		sectionData := p.data[sectionName]
		sortedKeys := make([]string, 0, len(sectionData))
		for k := range sectionData {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		fmt.Fprintf(&builder, "[%s]\n", sectionName)
		for _, key := range sortedKeys {
			fmt.Fprintf(&builder, "%s=%s\n", key, sectionData[key])
		}
	}

	return builder.String()
}

// LoadFromString loads the contents of a string into the parser and parses it to sections and keys-values.
func (p *Parser) LoadFromString(text string) error {
	input := bufio.NewScanner(strings.NewReader(text))
	err := p.parse(input)
	if err != nil {
		return fmt.Errorf("failed to parse input string: %w", err)
	}

	return nil
}

// ParseFile parses the given file in .ini format.
func (p *Parser) ParseFile(filePath string) error {
	if path.Ext(filePath) != ".ini" {
		return fmt.Errorf(".ini format is only support")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err = p.parse(scanner)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	return nil
}
