package graphiteparser

import (
	"fmt"
	"strconv"
	"strings"
)

// GraphiteParserError is used to report parser errors.
type GraphiteParserError struct {
	Message string
	Pos     int
}

func (e GraphiteParserError) Error() string {
	return fmt.Sprintf("Error at %d: %s", e.Pos, e.Message)
}

// AstNode represents a node in the abstract syntax tree.
type AstNode struct {
	Type        string      `json:"type"`
	Name        string      `json:"name,omitempty"`
	Params      []*AstNode  `json:"params,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Segments    []*AstNode  `json:"segments,omitempty"`
	Message     string      `json:"message,omitempty"`
	Pos         int         `json:"pos,omitempty"`
	IsUnclosed  bool        `json:"isUnclosed,omitempty"`
	Quote       string      `json:"quote,omitempty"`
	Base        int         `json:"base,omitempty"`
	IsMalformed bool        `json:"isMalformed,omitempty"`
}

// Parser holds the state of the parser.
type Parser struct {
	Expression string
	Lexer      *Lexer
	Tokens     []AstNode
	Index      int
}

// NewParser creates a new Parser instance.
func NewParser(expression string) *Parser {
	lexer := NewLexer(expression)
	tokens := lexer.Tokenize()
	return &Parser{
		Expression: expression,
		Lexer:      lexer,
		Tokens:     tokens,
		Index:      0,
	}
}

// GetAst returns the AST of the parsed expression.
// It recovers from parsing errors and returns an error node.
func (p *Parser) GetAst() (node *AstNode) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(GraphiteParserError); ok {
				node = &AstNode{
					Type:    "error",
					Message: err.Message,
					Pos:     err.Pos,
				}
			} else {
				// Unexpected panic; re-panic.
				panic(r)
			}
		}
	}()
	return p.start()
}

// start is the entry point for the parser.
// TODO: handle error
func (p *Parser) start() *AstNode {
	if node := p.functionCall(); node != nil {
		return node
	}
	if node := p.metricExpression(); node != nil {
		return node
	}
	return nil
}

// curlyBraceSegment parses a metric segment enclosed in curly braces.
func (p *Parser) curlyBraceSegment() *AstNode {
	if p.match("identifier", "{") || p.match("{") {
		curlySegment := ""

		// Loop until an empty token or a closing brace is encountered.
		for !p.match("") && !p.match("}") {
			token := p.consumeToken()
			curlySegment += token.Value.(string)
		}

		if !p.match("}") {
			p.errorMark("Expected closing '}'")
		}
		curlySegment += p.consumeToken().Value.(string)

		// If the curly segment is directly followed by an identifier, include it.
		if p.match("identifier") {
			curlySegment += p.consumeToken().Value.(string)
		}

		return &AstNode{
			Type:  "segment",
			Value: curlySegment,
		}
	}
	return nil
}

// metricSegment parses a metric segment.
func (p *Parser) metricSegment() *AstNode {
	if node := p.curlyBraceSegment(); node != nil {
		return node
	}

	if p.match("identifier") || p.match("number") || p.match("bool") {
		token := p.consumeToken()
		tokenValue, ok := token.Value.(string)
		var parts []string
		if !ok {
			parts = []string{""}
		} else {
			parts = strings.Split(tokenValue, ".")
		}
		if len(parts) == 2 {
			// Insert a dot token and a number token at the current index.
			dotToken := AstNode{Type: ".", Value: ".", Pos: token.Pos}
			numberToken := AstNode{Type: "number", Value: parts[1], Pos: token.Pos}
			// Insert the tokens into the slice.
			p.Tokens = append(p.Tokens[:p.Index],
				append([]AstNode{dotToken, numberToken}, p.Tokens[p.Index:]...)...,
			)
		}
		return &AstNode{
			Type:  "segment",
			Value: parts[0],
		}
	}

	if !p.match("templateStart") {
		p.errorMark("Expected metric identifier")
	}

	// Consume templateStart.
	p.consumeToken()

	if !p.match("identifier") {
		p.errorMark("Expected identifier after templateStart")
	}

	node := &AstNode{
		Type:  "template",
		Value: p.consumeToken().Value,
	}

	if !p.match("templateEnd") {
		p.errorMark("Expected templateEnd")
	}
	p.consumeToken() // Consume templateEnd.
	return node
}

// metricExpression parses a full metric expression.
func (p *Parser) metricExpression() *AstNode {
	if !p.match("templateStart") && !p.match("identifier") && !p.match("number") && !p.match("{") {
		return nil
	}

	node := &AstNode{
		Type:     "metric",
		Segments: []*AstNode{},
	}

	segment := p.metricSegment()
	if segment != nil {
		node.Segments = append(node.Segments, segment)
	}

	for p.match(".") {
		p.consumeToken() // Consume dot.
		segment := p.metricSegment()
		if segment == nil {
			p.errorMark("Expected metric identifier")
		}
		node.Segments = append(node.Segments, segment)
	}

	return node
}

// functionCall parses a function call.
func (p *Parser) functionCall() *AstNode {
	if !p.match("identifier", "(") {
		return nil
	}

	token := p.consumeToken()
	name := token.Value.(string)

	node := &AstNode{
		Type: "function",
		Name: name,
	}

	// Consume the left parenthesis.
	p.consumeToken()

	node.Params = p.functionParameters()

	if !p.match(")") {
		p.errorMark("Expected closing parenthesis")
	}
	p.consumeToken() // Consume the right parenthesis.
	return node
}

// boolExpression parses a boolean literal.
func (p *Parser) boolExpression() *AstNode {
	if !p.match("bool") {
		return nil
	}

	token := p.consumeToken()
	return &AstNode{
		Type:  "bool",
		Value: token.Value == "true",
	}
}

// functionParameters parses a comma-separated list of function parameters.
func (p *Parser) functionParameters() []*AstNode {
	if p.match(")") || p.match("") {
		return []*AstNode{}
	}

	var param *AstNode

	// Try the various possible parameter types.
	if node := p.functionCall(); node != nil {
		param = node
	} else if node := p.numericLiteral(); node != nil {
		param = node
	} else if node := p.seriesRefExpression(); node != nil {
		param = node
	} else if node := p.boolExpression(); node != nil {
		param = node
	} else if node := p.metricExpression(); node != nil {
		param = node
	} else if node := p.stringLiteral(); node != nil {
		param = node
	}

	// If no comma follows and a parameter was found, return it.
	if !p.match(",") && param != nil {
		return []*AstNode{param}
	}

	// Consume the comma.
	p.consumeToken()

	if param != nil {
		return append([]*AstNode{param}, p.functionParameters()...)
	}
	return []*AstNode{}
}

// seriesRefExpression parses a series reference.
func (p *Parser) seriesRefExpression() *AstNode {
	if !p.match("identifier") {
		return nil
	}

	// Peek at the token value.
	value := p.Tokens[p.Index].Value.(string)
	// Check if the value starts with '#' followed by an uppercase letter.
	if len(value) < 2 || !(strings.HasPrefix(value, "#") && strings.ToUpper(string(value[1])) == string(value[1])) {
		return nil
	}

	token := p.consumeToken()
	return &AstNode{
		Type:  "series-ref",
		Value: token.Value,
	}
}

// numericLiteral parses a numeric literal.
func (p *Parser) numericLiteral() *AstNode {
	if !p.match("number") {
		return nil
	}

	token := p.consumeToken()
	if token.Value != "" {
		if num, err := strconv.ParseFloat(token.Value.(string), 64); err == nil {
			return &AstNode{
				Type:  "number",
				Value: num,
			}
		}
	}
	return nil
}

// stringLiteral parses a string literal.
func (p *Parser) stringLiteral() *AstNode {
	if !p.match("string") {
		return nil
	}

	token := p.consumeToken()
	if token.IsUnclosed && token.Pos != 0 {
		p.errorMark("Unclosed string parameter")
	}
	return &AstNode{
		Type:  "string",
		Value: token.Value,
	}
}

// errorMark raises a parser error.
func (p *Parser) errorMark(text string) {
	var currentToken AstNode
	if p.Index < len(p.Tokens) {
		currentToken = p.Tokens[p.Index]
	} else {
		currentToken = AstNode{Type: "end of string", Pos: p.Lexer.Char}
	}
	err := GraphiteParserError{
		Message: text + " instead found " + currentToken.Type,
		Pos:     currentToken.Pos,
	}
	panic(err)
}

// consumeToken returns the current token and advances the parser.
func (p *Parser) consumeToken() AstNode {
	token := p.Tokens[p.Index]
	p.Index++
	return token
}

// matchToken checks whether the token at (current index + offset) matches the expected type.
func (p *Parser) matchToken(tokenType string, offset int) bool {
	index := p.Index + offset
	if index < 0 || index >= len(p.Tokens) {
		return false
	}
	t := p.Tokens[index].Type
	return t == tokenType
}

// match checks if the current token (and optionally the next) match the expected type(s).
// If a second token type is provided, it checks that the next token matches it.
func (p *Parser) match(token1 string, token2 ...string) bool {
	return p.matchToken(token1, 0) && (len(token2) == 0 || p.matchToken(token2[0], 1))
}
