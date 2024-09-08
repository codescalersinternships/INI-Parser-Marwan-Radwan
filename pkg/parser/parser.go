package parser

import (
	"bufio"
	"strings"
)

type Parser struct {
	data map[string]map[string]string
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		data: make(map[string]map[string]string),
	}
}

func (p *Parser) parse(scanner *bufio.Scanner) (map[string]map[string]string, error) {
	iniData := make(map[string]map[string]string)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines, comments and invalid lines
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse sections
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			iniData[currentSection] = make(map[string]string)
		} else if currentSection != "" {
			// Parse key=value pairs
			if equalIndex := strings.Index(line, "="); equalIndex != -1 {
				key := strings.TrimSpace(line[:equalIndex])
				value := strings.TrimSpace(line[equalIndex+1:])
				iniData[currentSection][key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return iniData, nil
}
