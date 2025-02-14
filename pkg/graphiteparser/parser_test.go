package graphiteparser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		expression string
		expected   *Parser
	}{
		{
			expression: "alias(movingAverage(statsd.fakesite.timers.ads_timer.upper_50, 5), 'OK')",
			expected: &Parser{
				Expression: "alias(movingAverage(statsd.fakesite.timers.ads_timer.upper_50, 5), 'OK')",
			},
		},
	}

	for _, test := range tests {
		parser := NewParser(test.expression)
		if parser.Expression != test.expected.Expression {
			t.Errorf("expected expression %s, got %s", test.expected.Expression, parser.Expression)
		}
	}
}
func TestGetAst(t *testing.T) {
	tests := []struct {
		expression string
		expected   *AstNode
	}{
		{
			expression: "alias(movingAverage(statsd.fakesite.timers.ads_timer.upper_50, 5), 'OK')",
			expected: &AstNode{
				Type: "function",
				Name: "alias",
				Params: []*AstNode{
					{
						Type: "function",
						Name: "movingAverage",
						Params: []*AstNode{
							{
								Type: "metric",
								Segments: []*AstNode{
									{Type: "segment", Value: "statsd"},
									{Type: "segment", Value: "fakesite"},
									{Type: "segment", Value: "timers"},
									{Type: "segment", Value: "ads_timer"},
									{Type: "segment", Value: "upper_50"},
								},
							},
							{Type: "number", Value: 5.0},
						},
					},
					{Type: "string", Value: "OK"},
				},
			},
		},
		{
			expression: "invalid(expression",
			expected: &AstNode{
				Type:    "error",
				Message: "Expected closing parenthesis instead found end of string",
				Pos:     19,
			},
		},
		{
			expression: "aliasByNode(apps.fakesite.web_server_01.counters.request_status.{code_500,code_400}.count, 5)",
			expected: &AstNode{
				Type: "function",
				Name: "aliasByNode",
				Params: []*AstNode{
					{
						Type: "metric",
						Segments: []*AstNode{
							{Type: "segment", Value: "apps"},
							{Type: "segment", Value: "fakesite"},
							{Type: "segment", Value: "web_server_01"},
							{Type: "segment", Value: "counters"},
							{Type: "segment", Value: "request_status"},
							{Type: "segment", Value: "{code_500,code_400}"},
							{Type: "segment", Value: "count"},
						},
					},
					{Type: "number", Value: 5.0},
				},
			},
		},
	}

	for _, test := range tests {
		parser := NewParser(test.expression)
		ast := parser.GetAst()
		if !compareAstNodes(ast, test.expected) {
			t.Errorf("expected AST %v, got %v", test.expected, ast)
		}
	}
}

// compareAstNodes is a helper function to compare two AST nodes.
func compareAstNodes(a, b *AstNode) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Type != b.Type || a.Name != b.Name || a.Value != b.Value || a.Message != b.Message || a.Pos != b.Pos {
		return false
	}
	if len(a.Params) != len(b.Params) || len(a.Segments) != len(b.Segments) {
		return false
	}
	for i := range a.Params {
		if !compareAstNodes(a.Params[i], b.Params[i]) {
			return false
		}
	}
	for i := range a.Segments {
		if !compareAstNodes(a.Segments[i], b.Segments[i]) {
			return false
		}
	}
	return true
}
