package parser

import (
	"bufio"
	"fmt"
	"strings"
)

type Parser struct {
	data       map[string]map[string]string
	globalKeys map[string]string
	sections   []string
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
			if _, exist := p.data[currentSection]; !exist {
				p.data[currentSection] = make(map[string]string)
				p.sections = append(p.sections, currentSection)
			}
		} else {
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
		}
	}

	return scanner.Err()
}
