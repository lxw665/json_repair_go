// Package jsonrepair provides functionality to repair invalid JSON strings
// This module will parse the JSON file following the BNF definition:
//
//	<json> ::= <primitive> | <container>
//	<primitive> ::= <number> | <string> | <boolean>
//	<container> ::= <object> | <array>
//	<array> ::= '[' [ <json> *(', ' <json>) ] ']'
//	<object> ::= '{' [ <member> *(', ' <member>) ] '}'
//	<member> ::= <string> ': ' <json>
//
// If something is wrong (missing parentheses or quotes) it will use simple heuristics to fix the JSON string:
// - Add missing parentheses if the parser believes array or object should be closed
// - Quote strings or add missing quotes
// - Adjust whitespaces and remove line breaks
package jsonrepair

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RepairOptions contains options for JSON repair operations
type RepairOptions struct {
	SkipValidation bool // Skip validation of input as valid JSON
	ReturnObjects  bool // Return parsed objects instead of JSON strings
	Logging        bool // Enable detailed logging of repair operations
	StreamStable   bool // Keep repair results stable for streaming JSON
	EnsureASCII    bool // Ensure non-ASCII characters are escaped
}

// RepairJSON repairs a malformed JSON string and returns the corrected JSON
func RepairJSON(jsonStr string) (string, error) {
	return RepairJSONWithOptions(jsonStr, RepairOptions{})
}

// RepairJSONWithOptions repairs a malformed JSON string with specified options
func RepairJSONWithOptions(jsonStr string, options RepairOptions) (string, error) {
	parser := NewJSONParser(jsonStr, options)

	// Try to parse as valid JSON first if not skipping validation
	if !options.SkipValidation {
		var temp interface{}
		if err := json.Unmarshal([]byte(jsonStr), &temp); err == nil {
			// Already valid JSON
			if options.ReturnObjects {
				return jsonStr, nil
			}
			// Re-marshal to ensure consistent formatting
			result, err := json.Marshal(temp)
			return string(result), err
		}
	}

	// Parse with repair
	result, err := parser.Parse()
	if err != nil {
		return "", err
	}

	// Return as JSON string
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal repaired result: %w", err)
	}

	return string(jsonBytes), nil
}

// Loads parses a JSON string (with repair if needed) and returns the parsed object
func Loads(jsonStr string) (interface{}, error) {
	parser := NewJSONParser(jsonStr, RepairOptions{ReturnObjects: true})
	return parser.Parse()
}

// Load reads JSON from an io.Reader and parses it (with repair if needed) into v
func Load(reader io.Reader, v interface{}) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return Unmarshal(data, v)
}

// Unmarshal parses JSON data (with repair if needed) and stores the result in v
func Unmarshal(data []byte, v interface{}) error {
	repaired, err := RepairJSON(string(data))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(repaired), v)
}

// FromFile reads and repairs JSON from a file
func FromFile(filename string) (interface{}, error) {
	// This would be implemented to read from file
	// For now, return an error indicating it's not implemented
	return nil, fmt.Errorf("FromFile not implemented yet")
}

// JSONParser handles parsing and repairing JSON strings
type JSONParser struct {
	jsonStr string
	index   int
	options RepairOptions
	logs    []map[string]string
	context string
}

// NewJSONParser creates a new JSON parser with the given input and options
func NewJSONParser(jsonStr string, options RepairOptions) *JSONParser {
	return &JSONParser{
		jsonStr: jsonStr,
		index:   0,
		options: options,
		logs:    make([]map[string]string, 0),
	}
}

// Parse parses the JSON string and returns the parsed result
func (p *JSONParser) Parse() (interface{}, error) {
	result, err := p.parseJSON()
	if err != nil {
		return nil, err
	}

	// Check if there are multiple JSON objects/values
	if p.index < len(p.jsonStr) {
		p.log("The parser returned early, checking if there's more json elements")
		results := []interface{}{result}

		for p.index < len(p.jsonStr) {
			nextResult, err := p.parseJSON()
			if err != nil {
				p.index++
				continue
			}
			if nextResult != nil && nextResult != "" {
				results = append(results, nextResult)
			}
		}

		if len(results) == 1 {
			return results[0], nil
		}
		return results, nil
	}

	return result, nil
}

// parseJSON parses a single JSON value
func (p *JSONParser) parseJSON() (interface{}, error) {
	for p.index < len(p.jsonStr) {
		char := p.getCharAt()
		if char == 0 {
			return "", nil
		}

		switch char {
		case '{':
			p.index++
			return p.parseObject()
		case '[':
			p.index++
			return p.parseArray()
		case '"', '\'':
			return p.parseString()
		case '#', '/':
			p.parseComment()
			continue
		default:
			if p.isDigit(char) || char == '-' || char == '.' {
				return p.parseNumber()
			} else if p.isAlpha(char) {
				return p.parseString()
			} else {
				p.index++
			}
		}
	}
	return "", nil
}

// parseObject parses a JSON object
func (p *JSONParser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	p.skipWhitespaces()

	// Handle empty object
	if p.getCharAt() == '}' {
		p.index++
		return obj, nil
	}

	for {
		p.skipWhitespaces()

		// Check for end of object
		char := p.getCharAt()
		if char == 0 || char == '}' {
			if char == '}' {
				p.index++
			}
			break
		}

		// Parse key
		key, err := p.parseString()
		if err != nil {
			break
		}
		keyStr, ok := key.(string)
		if !ok {
			keyStr = fmt.Sprintf("%v", key)
		}

		p.skipWhitespaces()

		// Expect colon
		if p.getCharAt() != ':' {
			p.log("While parsing object, expected ':' after key")
			// Add missing colon
		} else {
			p.index++
		}

		p.skipWhitespaces()

		// Parse value
		value, err := p.parseJSON()
		if err != nil {
			value = ""
		}

		obj[keyStr] = value

		p.skipWhitespaces()

		// Check for comma or end
		char = p.getCharAt()
		if char == ',' {
			p.index++
			continue
		} else if char == '}' {
			p.index++
			break
		} else if char == 0 {
			break
		}
		// Missing comma, continue anyway
	}

	return obj, nil
}

// parseArray parses a JSON array
func (p *JSONParser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)
	p.skipWhitespaces()

	// Handle empty array
	if p.getCharAt() == ']' {
		p.index++
		return arr, nil
	}

	for {
		p.skipWhitespaces()

		char := p.getCharAt()
		if char == 0 || char == ']' {
			if char == ']' {
				p.index++
			} else {
				p.log("While parsing an array we missed the closing ], ignoring it")
			}
			break
		}

		// Parse value
		value, err := p.parseJSON()
		if err != nil {
			p.index++
			continue
		}

		if value == "" {
			p.index++
			continue
		}

		arr = append(arr, value)

		p.skipWhitespaces()

		// Check for comma or end
		char = p.getCharAt()
		if char == ',' {
			p.index++
			continue
		} else if char == ']' {
			p.index++
			break
		} else if char == 0 {
			break
		}
		// Missing comma, continue anyway
	}

	return arr, nil
}

// parseString parses a JSON string (quoted or unquoted)
func (p *JSONParser) parseString() (interface{}, error) {
	char := p.getCharAt()

	// Check for boolean values
	if p.matchWord("true") {
		p.index += 4
		return true, nil
	}
	if p.matchWord("false") {
		p.index += 5
		return false, nil
	}
	if p.matchWord("True") {
		p.index += 4
		return true, nil
	}
	if p.matchWord("False") {
		p.index += 5
		return false, nil
	}
	if p.matchWord("null") {
		p.index += 4
		return nil, nil
	}
	if p.matchWord("None") {
		p.index += 4
		return nil, nil
	}

	// Handle quoted strings
	if char == '"' || char == '\'' {
		return p.parseQuotedString(char)
	}

	// Handle unquoted strings
	return p.parseUnquotedString()
}

// parseQuotedString parses a quoted string
func (p *JSONParser) parseQuotedString(quote byte) (string, error) {
	p.index++ // skip opening quote
	var result strings.Builder

	for p.index < len(p.jsonStr) {
		char := p.jsonStr[p.index]

		if char == quote {
			p.index++ // skip closing quote
			return result.String(), nil
		}

		if char == '\\' && p.index+1 < len(p.jsonStr) {
			p.index++
			next := p.jsonStr[p.index]
			switch next {
			case '"', '\'', '\\', '/':
				result.WriteByte(next)
			case 'b':
				result.WriteByte('\b')
			case 'f':
				result.WriteByte('\f')
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case 'u':
				// Unicode escape
				if p.index+4 < len(p.jsonStr) {
					hex := p.jsonStr[p.index+1 : p.index+5]
					if code, err := strconv.ParseInt(hex, 16, 32); err == nil {
						result.WriteRune(rune(code))
						p.index += 4
					} else {
						result.WriteByte(next)
					}
				} else {
					result.WriteByte(next)
				}
			default:
				result.WriteByte(next)
			}
		} else {
			result.WriteByte(char)
		}
		p.index++
	}

	p.log("While parsing a string, we missed the closing quote, ignoring")
	return result.String(), nil
}

// parseUnquotedString parses an unquoted string
func (p *JSONParser) parseUnquotedString() (string, error) {
	var result strings.Builder
	start := p.index

	for p.index < len(p.jsonStr) {
		char := p.getCharAt()
		if char == ',' || char == '}' || char == ']' || char == ':' || p.isWhitespace(char) {
			break
		}
		result.WriteByte(char)
		p.index++
	}

	if p.index == start {
		return "", nil
	}

	return result.String(), nil
}

// parseNumber parses a number value
func (p *JSONParser) parseNumber() (interface{}, error) {
	start := p.index

	// Handle sign
	if p.getCharAt() == '-' {
		p.index++
	}

	// Parse integer part
	hasDigits := false
	for p.index < len(p.jsonStr) && p.isDigit(p.getCharAt()) {
		hasDigits = true
		p.index++
	}

	if !hasDigits {
		p.index = start
		return p.parseUnquotedString()
	}

	// Check for decimal
	hasDecimal := false
	if p.getCharAt() == '.' {
		hasDecimal = true
		p.index++
		for p.index < len(p.jsonStr) && p.isDigit(p.getCharAt()) {
			p.index++
		}
	}

	// Check for exponent
	char := p.getCharAt()
	if char == 'e' || char == 'E' {
		p.index++
		char = p.getCharAt()
		if char == '+' || char == '-' {
			p.index++
		}
		for p.index < len(p.jsonStr) && p.isDigit(p.getCharAt()) {
			p.index++
		}
	}

	numStr := p.jsonStr[start:p.index]

	if hasDecimal {
		if val, err := strconv.ParseFloat(numStr, 64); err == nil {
			return val, nil
		}
	} else {
		if val, err := strconv.ParseInt(numStr, 10, 64); err == nil {
			return val, nil
		}
	}

	return numStr, nil
}

// parseComment skips over comments
func (p *JSONParser) parseComment() {
	char := p.getCharAt()

	if char == '/' && p.index+1 < len(p.jsonStr) {
		next := p.jsonStr[p.index+1]
		if next == '/' {
			// Line comment
			p.index += 2
			for p.index < len(p.jsonStr) && p.getCharAt() != '\n' {
				p.index++
			}
		} else if next == '*' {
			// Block comment
			p.index += 2
			for p.index+1 < len(p.jsonStr) {
				if p.getCharAt() == '*' && p.jsonStr[p.index+1] == '/' {
					p.index += 2
					break
				}
				p.index++
			}
		}
	} else if char == '#' {
		// Python-style comment
		for p.index < len(p.jsonStr) && p.getCharAt() != '\n' {
			p.index++
		}
	}
}

// Helper methods

// getCharAt returns the character at the current index, or 0 if at end
func (p *JSONParser) getCharAt() byte {
	if p.index >= len(p.jsonStr) {
		return 0
	}
	return p.jsonStr[p.index]
}

// skipWhitespaces skips whitespace characters
func (p *JSONParser) skipWhitespaces() {
	for p.index < len(p.jsonStr) && p.isWhitespace(p.getCharAt()) {
		p.index++
	}
}

// isWhitespace checks if a character is whitespace
func (p *JSONParser) isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

// isDigit checks if a character is a digit
func (p *JSONParser) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

// isAlpha checks if a character is alphabetic
func (p *JSONParser) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

// matchWord checks if the current position matches a word
func (p *JSONParser) matchWord(word string) bool {
	if p.index+len(word) > len(p.jsonStr) {
		return false
	}
	return strings.EqualFold(p.jsonStr[p.index:p.index+len(word)], word)
}

// log adds a log entry if logging is enabled
func (p *JSONParser) log(message string) {
	if p.options.Logging {
		logEntry := map[string]string{
			"text":    message,
			"context": p.jsonStr,
		}
		p.logs = append(p.logs, logEntry)
	}
}
