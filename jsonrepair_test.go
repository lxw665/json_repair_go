package jsonrepair_test

import (
"encoding/json"
"reflect"
"testing"

"github.com/lxw/json_repair_go"
)

// TestBasicRepairs tests basic JSON repair functionality
func TestBasicRepairs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		// Basic valid JSON should remain unchanged
		{
			name:     "valid_json",
			input:    `{"a": 1}`,
			expected: map[string]interface{}{"a": int64(1)},
		},
		
		// Missing closing brace
		{
			name:     "missing_closing_brace",
			input:    `{"a": 1`,
			expected: map[string]interface{}{"a": int64(1)},
		},
		
		// Missing quotes around keys
		{
			name:     "unquoted_keys",
			input:    `{a: 1}`,
			expected: map[string]interface{}{"a": int64(1)},
		},
		
		// Single quotes instead of double quotes
		{
			name:     "single_quotes",
			input:    `{'a': 1}`,
			expected: map[string]interface{}{"a": int64(1)},
		},
		
		// Boolean values
		{
			name:     "boolean_True",
			input:    `{"active": True}`,
			expected: map[string]interface{}{"active": true},
		},
		
		{
			name:     "boolean_False",
			input:    `{"active": False}`,
			expected: map[string]interface{}{"active": false},
		},
		
		// Null values
		{
			name:     "None_value",
			input:    `{"value": None}`,
			expected: map[string]interface{}{"value": nil},
		},
		
		// String values
		{
			name:     "unquoted_string",
			input:    `{"name": John}`,
			expected: map[string]interface{}{"name": "John"},
		},

		// Float values
		{
			name:     "float_value",
			input:    `{"price": 19.99}`,
			expected: map[string]interface{}{"price": 19.99},
		},

		// Arrays
		{
			name:     "array_missing_bracket",
			input:    `[1, 2, 3`,
			expected: []interface{}{int64(1), int64(2), int64(3)},
		},

		// Complex cases
		{
			name:     "mixed_issues",
			input:    `{name: 'John', age: 30, active: True`,
			expected: map[string]interface{}{"name": "John", "age": int64(30), "active": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
result, err := jsonrepair.Loads(tt.input)
if err != nil {
t.Errorf("Loads() error = %v", err)
return
}

if !reflect.DeepEqual(result, tt.expected) {
t.Errorf("Loads() = %v (type: %T), want %v (type: %T)", result, result, tt.expected, tt.expected)
}
})
	}
}

// TestRepairJSON tests the RepairJSON function
func TestRepairJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid_json",
			input: `{"a": 1}`,
		},
		{
			name:  "missing_closing_brace",
			input: `{"a": 1`,
		},
		{
			name:  "unquoted_keys",
			input: `{a: 1}`,
		},
		{
			name:  "single_quotes",
			input: `{'a': 1}`,
		},
		{
			name:  "boolean_values",
			input: `{"active": True, "disabled": False}`,
		},
		{
			name:  "null_value",
			input: `{"value": None}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
result, err := jsonrepair.RepairJSON(tt.input)
if err != nil {
t.Errorf("RepairJSON() error = %v", err)
return
}

// Verify the result is valid JSON
var parsed interface{}
if err := json.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("RepairJSON() produced invalid JSON: %s, error: %v", result, err)
				return
			}

			t.Logf("Input: %s -> Output: %s", tt.input, result)
		})
	}
}

// TestErrorCases tests error handling
func TestErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		shouldErr bool
	}{
		{
			name:      "empty_string",
			input:     "",
			shouldErr: false, // Should handle gracefully
		},
		{
			name:      "whitespace_only",
			input:     "   \n\t  ",
			shouldErr: false,
		},
		{
			name:      "just_text",
			input:     "hello world",
			shouldErr: false, // Should parse as string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
result, err := jsonrepair.Loads(tt.input)

if tt.shouldErr {
if err == nil {
t.Errorf("Expected error but got none")
}
} else {
if err != nil {
t.Errorf("Unexpected error: %v", err)
} else {
t.Logf("Input: '%s' -> Result: %v", tt.input, result)
}
}
})
	}
}
