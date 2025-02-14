package graphiteparser

import (
	"math"
	"testing"
)

func TestIsDecimalDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"5", true},
		{"9", true},
		{"a", false},
		{"", false},
		{"10", false},
		{" ", false},
		{"$", false},
	}

	for _, test := range tests {
		result := isDecimalDigit(test.input)
		if result != test.expected {
			t.Errorf("isDecimalDigit(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestIsOctalDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"3", true},
		{"7", true},
		{"8", false},
		{"9", false},
		{"a", false},
		{"", false},
		{"10", false},
		{" ", false},
		{"$", false},
	}

	for _, test := range tests {
		result := isOctalDigit(test.input)
		if result != test.expected {
			t.Errorf("isOctalDigit(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestIsHexDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"9", true},
		{"a", true},
		{"f", true},
		{"A", true},
		{"F", true},
		{"g", false},
		{"G", false},
		{"", false},
		{"10", false},
		{" ", false},
		{"$", false},
	}

	for _, test := range tests {
		result := isHexDigit(test.input)
		if result != test.expected {
			t.Errorf("isHexDigit(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestIsIdentifierStart(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"$", true},
		{"_", true},
		{"\\", true},
		{"a", true},
		{"z", true},
		{"A", true},
		{"Z", true},
		{"0", false},
		{"9", false},
		{"-", false},
		{" ", false},
		{"", false},
		{"1a", false},
	}

	for _, test := range tests {
		result := isIdentifierStart(test.input)
		if result != test.expected {
			t.Errorf("isIdentifierStart(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestIsPunctuator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{".", true},
		{"(", true},
		{")", true},
		{",", true},
		{"{", true},
		{"}", true},
		{";", false},
		{"a", false},
		{"1", false},
		{"", false},
		{"[]", false},
	}

	for _, test := range tests {
		result := isPunctuator(test.input)
		if result != test.expected {
			t.Errorf("isPunctuator(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestIsFinite(t *testing.T) {
	tests := []struct {
		input    float64
		expected bool
	}{
		{0.0, true},
		{1.0, true},
		{-1.0, true},
		{math.Inf(1), false},
		{math.Inf(-1), false},
		{math.NaN(), false},
		{math.MaxFloat64, true},
		{-math.MaxFloat64, true},
	}

	for _, test := range tests {
		result := isFinite(test.input)
		if result != test.expected {
			t.Errorf("isFinite(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
