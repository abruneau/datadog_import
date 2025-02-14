package graphiteparser

import (
	"regexp"
	"strconv"
)

// -----------------------------------------------------------------------------
// Unicode tables and identifier tables (auto-generated from Unicode data)
// -----------------------------------------------------------------------------

// unicodeLetterTable contains paired ranges (start, end) of Unicode code points
// for letters (e.g. Lu, Ll, Lt, Lm, Lo, Nl).
var unicodeLetterTable = []int{
	170, 170, 181, 181, 186, 186, 192, 214, 216, 246, 248, 705, 710, 721, 736, 740, 748, 748, 750, 750, 880, 884, 886,
	887, 890, 893, 902, 902, 904, 906, 908, 908, 910, 929, 931, 1013, 1015, 1153, 1162, 1319, 1329, 1366, 1369, 1369,
	1377, 1415, 1488, 1514, 1520, 1522, 1568, 1610, 1646, 1647, 1649, 1747, 1749, 1749, 1765, 1766, 1774, 1775, 1786,
	1788, 1791, 1791, 1808, 1808, 1810, 1839, 1869, 1957, 1969, 1969, 1994, 2026, 2036, 2037, 2042, 2042, 2048, 2069,
	2074, 2074, 2084, 2084, 2088, 2088, 2112, 2136, 2308, 2361, 2365, 2365, 2384, 2384, 2392, 2401, 2417, 2423, 2425,
	2431, 2437, 2444, 2447, 2448, 2451, 2472, 2474, 2480, 2482, 2482, 2486, 2489, 2493, 2493, 2510, 2510, 2524, 2525,
	2527, 2529, 2544, 2545, 2565, 2570, 2575, 2576, 2579, 2600, 2602, 2608, 2610, 2611, 2613, 2614, 2616, 2617, 2649,
	2652, 2654, 2654, 2674, 2676, 2693, 2701, 2703, 2705, 2707, 2728, 2730, 2736, 2738, 2739, 2741, 2745, 2749, 2749,
	2768, 2768, 2784, 2785, 2821, 2828, 2831, 2832, 2835, 2856, 2858, 2864, 2866, 2867, 2869, 2873, 2877, 2877, 2908,
	2909, 2911, 2913, 2929, 2929, 2947, 2947, 2949, 2954, 2958, 2960, 2962, 2965, 2969, 2970, 2972, 2972, 2974, 2975,
	2979, 2980, 2984, 2986, 2990, 3001, 3024, 3024, 3077, 3084, 3086, 3088, 3090, 3112, 3114, 3123, 3125, 3129, 3133,
	3133, 3160, 3161, 3168, 3169, 3205, 3212, 3214, 3216, 3218, 3240, 3242, 3251, 3253, 3257, 3261, 3261, 3294, 3294,
	3296, 3297, 3313, 3314, 3333, 3340, 3342, 3344, 3346, 3386, 3389, 3389, 3406, 3406, 3424, 3425, 3450, 3455, 3461,
	3478, 3482, 3505, 3507, 3515, 3517, 3517, 3520, 3526, 3585, 3632, 3634, 3635, 3648, 3654, 3713, 3714, 3716, 3716,
	3719, 3720, 3722, 3722, 3725, 3725, 3732, 3735, 3737, 3743, 3745, 3747, 3749, 3749, 3751, 3751, 3754, 3755, 3757,
	3760, 3762, 3763, 3773, 3773, 3776, 3780, 3782, 3782, 3804, 3805, 3840, 3840, 3904, 3911, 3913, 3948, 3976, 3980,
	4096, 4138, 4159, 4159, 4176, 4181, 4186, 4189, 4193, 4193, 4197, 4198, 4206, 4208, 4213, 4225, 4238, 4238, 4256,
	4293, 4304, 4346, 4348, 4348, 4352, 4680, 4682, 4685, 4688, 4694, 4696, 4696, 4698, 4701, 4704, 4744, 4746, 4749,
	4752, 4784, 4786, 4789, 4792, 4798, 4800, 4800, 4802, 4805, 4808, 4822, 4824, 4880, 4882, 4885, 4888, 4954, 4992,
	5007, 5024, 5108, 5121, 5740, 5743, 5759, 5761, 5786, 5792, 5866, 5870, 5872, 5888, 5900, 5902, 5905, 5920, 5937,
	5952, 5969, 5984, 5996, 5998, 6000, 6016, 6067, 6103, 6103, 6108, 6108, 6176, 6263, 6272, 6312, 6314, 6314, 6320,
	6389, 6400, 6428, 6480, 6509, 6512, 6516, 6528, 6571, 6593, 6599, 6656, 6678, 6688, 6740, 6823, 6823, 6917, 6963,
	6981, 6987, 7043, 7072, 7086, 7087, 7104, 7141, 7168, 7203, 7245, 7247, 7258, 7293, 7401, 7404, 7406, 7409, 7424,
	7615, 7680, 7957, 7960, 7965, 7968, 8005, 8008, 8013, 8016, 8023, 8025, 8025, 8027, 8027, 8029, 8029, 8031, 8061,
	8064, 8116, 8118, 8124, 8126, 8126, 8130, 8132, 8134, 8140, 8144, 8147, 8150, 8155, 8160, 8172, 8178, 8180, 8182,
	8188, 8305, 8305, 8319, 8319, 8336, 8348, 8450, 8450, 8455, 8455, 8458, 8467, 8469, 8469, 8473, 8477, 8484, 8484,
	8486, 8486, 8488, 8488, 8490, 8493, 8495, 8505, 8508, 8511, 8517, 8521, 8526, 8526, 8544, 8584, 11264, 11310, 11312,
	11358, 11360, 11492, 11499, 11502, 11520, 11557, 11568, 11621, 11631, 11631, 11648, 11670, 11680, 11686, 11688, 11694,
	11696, 11702, 11704, 11710, 11712, 11718, 11720, 11726, 11728, 11734, 11736, 11742, 11823, 11823, 12293, 12295, 12321,
	12329, 12337, 12341, 12344, 12348, 12353, 12438, 12445, 12447, 12449, 12538, 12540, 12543, 12549, 12589, 12593, 12686,
	12704, 12730, 12784, 12799, 13312, 13312, 19893, 19893, 19968, 19968, 40907, 40907, 40960, 42124, 42192, 42237, 42240,
	42508, 42512, 42527, 42538, 42539, 42560, 42606, 42623, 42647, 42656, 42735, 42775, 42783, 42786, 42888, 42891, 42894,
	42896, 42897, 42912, 42921, 43002, 43009, 43011, 43013, 43015, 43018, 43020, 43042, 43072, 43123, 43138, 43187, 43250,
	43255, 43259, 43259, 43274, 43301, 43312, 43334, 43360, 43388, 43396, 43442, 43471, 43471, 43520, 43560, 43584, 43586,
	43588, 43595, 43616, 43638, 43642, 43642, 43648, 43695, 43697, 43697, 43701, 43702, 43705, 43709, 43712, 43712, 43714,
	43714, 43739, 43741, 43777, 43782, 43785, 43790, 43793, 43798, 43808, 43814, 43816, 43822, 43968, 44002, 44032, 44032,
	55203, 55203, 55216, 55238, 55243, 55291, 63744, 64045, 64048, 64109, 64112, 64217, 64256, 64262, 64275, 64279, 64285,
	64285, 64287, 64296, 64298, 64310, 64312, 64316, 64318, 64318, 64320, 64321, 64323, 64324, 64326, 64433, 64467, 64829,
	64848, 64911, 64914, 64967, 65008, 65019, 65136, 65140, 65142, 65276, 65313, 65338, 65345, 65370, 65382, 65470, 65474,
	65479, 65482, 65487, 65490, 65495, 65498, 65500, 65536, 65547, 65549, 65574, 65576, 65594, 65596, 65597, 65599, 65613,
	65616, 65629, 65664, 65786, 65856, 65908, 66176, 66204, 66208, 66256, 66304, 66334, 66352, 66378, 66432, 66461, 66464,
	66499, 66504, 66511, 66513, 66517, 66560, 66717, 67584, 67589, 67592, 67592, 67594, 67637, 67639, 67640, 67644, 67644,
	67647, 67669, 67840, 67861, 67872, 67897, 68096, 68096, 68112, 68115, 68117, 68119, 68121, 68147, 68192, 68220, 68352,
	68405, 68416, 68437, 68448, 68466, 68608, 68680, 69635, 69687, 69763, 69807, 73728, 74606, 74752, 74850, 77824, 78894,
	92160, 92728, 110592, 110593, 119808, 119892, 119894, 119964, 119966, 119967, 119970, 119970, 119973, 119974, 119977,
	119980, 119982, 119993, 119995, 119995, 119997, 120003, 120005, 120069, 120071, 120074, 120077, 120084, 120086,
	120092, 120094, 120121, 120123, 120126, 120128, 120132, 120134, 120134, 120138, 120144, 120146, 120485, 120488,
	120512, 120514, 120538, 120540, 120570, 120572, 120596, 120598, 120628, 120630, 120654, 120656, 120686, 120688,
	120712, 120714, 120744, 120746, 120770, 120772, 120779, 131072, 131072, 173782, 173782, 173824, 173824, 177972,
	177972, 177984, 177984, 178205, 178205, 194560, 195101,
}

// identifierStartTable for ASCII (0-127) characters.
// true indicates that the character can start an identifier.
var identifierStartTable [128]bool

func initTable() {
	for i := 0; i < 128; i++ {
		identifierStartTable[i] =
			(i >= 48 && i <= 57) || // 0-9
				i == 36 || // $
				i == 126 || // ~
				i == 124 || // |
				(i >= 65 && i <= 90) || // A-Z
				i == 95 || // _
				i == 45 || // -
				i == 42 || // *
				i == 58 || // :
				i == 91 || // templateStart [
				i == 93 || // templateEnd ]
				i == 63 || // ?
				i == 37 || // %
				i == 35 || // #
				i == 61 || // =
				i == 64 || // @
				(i >= 97 && i <= 122) // a-z
	}
}

// For this example, we use the same table for identifier parts.
var identifierPartTable [128]bool

// -----------------------------------------------------------------------------
// Token and Lexer types
// -----------------------------------------------------------------------------

// Lexer holds the input string and a current position.
type Lexer struct {
	input string
	Char  int // current character (1-indexed)
	from  int // start of current token
}

// NewLexer creates a new Lexer for the given expression.
func NewLexer(expression string) *Lexer {
	initTable()
	identifierPartTable = identifierStartTable
	return &Lexer{
		input: expression,
		Char:  1,
		from:  1,
	}
}

func (l *Lexer) peek(i int) string {
	if len(l.input) > i {
		return string(l.input[i])
	}
	return ""
}

// skip advances the lexer position by n characters (default 1).
func (l *Lexer) skip(n int) {
	l.Char += n
	l.input = l.input[n:]
}

// Tokenize produces a slice of tokens from the input.
func (l *Lexer) Tokenize() []AstNode {
	var tokens []AstNode
	for {
		tok := l.next()
		if tok == nil {
			break
		}
		tokens = append(tokens, *tok)
	}
	return tokens
}

// next returns the next token from the input (or nil if end-of-input).
func (l *Lexer) next() *AstNode {
	// Record the starting position (1-indexed for pos reporting).
	l.from = l.Char

	if l.peek(0) == "" {
		return nil
	}

	// Skip whitespace.
	whitespaceRegexp := regexp.MustCompile(`\s`)
	if whitespaceRegexp.MatchString(l.peek(0)) {
		for {
			if whitespaceRegexp.MatchString(l.peek(0)) {
				l.from += 1
				l.skip(1)
			} else {
				break
			}
		}
	}

	// Attempt to scan a string literal.
	if tok := l.scanStringLiteral(); tok != nil {
		return tok
	}

	// Try punctuator, numeric literal, identifier, or template sequence.
	if tok := l.scanPunctuator(); tok != nil {
		l.skip(len(tok.Value.(string)))
		return tok
	}
	if tok := l.scanNumericLiteral(); tok != nil {
		l.skip(len(tok.Value.(string)))
		return tok
	}
	if tok := l.scanIdentifier(); tok != nil {
		l.skip(len(tok.Value.(string)))
		return tok
	}
	if tok := l.scanTemplateSequence(); tok != nil {
		l.skip(len(tok.Value.(string)))
		return tok
	}

	// No token could be matched.
	return nil
}

// scanTemplateSequence returns a template token if one is found.
func (l *Lexer) scanTemplateSequence() *AstNode {
	if l.peek(0) == "[" && l.peek(1) == "[" {
		return &AstNode{
			Type:  "templateStart",
			Value: "[[",
			Pos:   l.Char,
		}
	}
	if l.peek(0) == "]" && l.peek(1) == "]" {
		// (Note: original code returns '[[' as value here; reproduced as-is.)
		return &AstNode{
			Type:  "templateEnd",
			Value: "[[",
			Pos:   l.Char,
		}
	}
	return nil
}

// scanIdentifier extracts an identifier (or boolean literal) token.
func (l *Lexer) scanIdentifier() *AstNode {
	// local index relative to current l.pos.
	index := 0
	var id string

	// Helper: check if code is a Unicode letter using unicodeLetterTable.
	isUnicodeLetter := func(code int) bool {
		for i := 0; i < len(unicodeLetterTable); {
			if code < unicodeLetterTable[i+1] {
				return false
			}
			if code <= unicodeLetterTable[i+1] {
				return true
			}
		}
		return false
	}

	// Helper: isHexDigit
	isHexDigit := func(s string) bool {
		matched, _ := regexp.MatchString(`^[0-9a-fA-F]$`, s)
		return matched
	}

	// readUnicodeEscapeSequence attempts to read a Unicode escape (e.g. "\u0041")
	readUnicodeEscapeSequence := func() string {
		// Expect a backslash already encountered.
		index++ // skip the backslash
		if l.peek(index) != "u" {
			return ""
		}
		// We expect 4 hex digits.
		if len(l.input) < l.Char+index+4 {
			return ""
		}
		ch1 := l.peek(index + 1)
		ch2 := l.peek(index + 2)
		ch3 := l.peek(index + 3)
		ch4 := l.peek(index + 4)

		if isHexDigit(ch1) && isHexDigit(ch2) && isHexDigit(ch3) && isHexDigit(ch4) {
			code, err := strconv.ParseInt(ch1+ch2+ch3+ch4, 16, 32)
			if err != nil {
				return ""
			}
			if isUnicodeLetter(int(code)) {
				index += 5
				return "\\u" + ch1 + ch2 + ch3 + ch4
			}
			return ""
		}
		return ""
	}

	// getIdentifierStart returns the first character (or escape sequence) if valid.
	getIdentifierStart := func() string {
		chr := l.peek(index)
		if chr == "*" {
			index++
			return chr
		}
		code := int(chr[0])
		if code == 92 { // backslash
			return readUnicodeEscapeSequence()
		}
		if code < 128 {
			if identifierStartTable[code] {
				index++
				return chr
			}
			return ""
		}
		// For non-ASCII, check if Unicode letter.
		if isUnicodeLetter(code) {
			index++
			return chr
		}
		return ""
	}

	// getIdentifierPart returns subsequent valid identifier characters.
	getIdentifierPart := func() string {
		chr := l.peek(index)
		if chr == "" {
			return ""
		}
		code := int(chr[0])
		if code == 92 {
			return readUnicodeEscapeSequence()
		}
		if code < 128 {
			if identifierPartTable[code] {
				index++
				return chr
			}
			return ""
		}
		if isUnicodeLetter(code) {
			index++
			return chr
		}
		return ""
	}

	// Get the first character.
	ch := getIdentifierStart()
	if ch == "" {
		return nil
	}

	id = ch

	// Consume subsequent identifier characters.
	for {
		ch := getIdentifierPart()
		if ch != "" {
			id += ch
		} else {
			break
		}
	}

	// Determine the token type.
	tokType := "identifier"
	if id == "true" || id == "false" {
		tokType = "bool"
	}

	return &AstNode{
		Type:  tokType,
		Value: id,
		Pos:   l.Char,
	}
}

// scanNumericLiteral extracts a numeric literal token.
func (l *Lexer) scanNumericLiteral() *AstNode {
	index := 0
	var value string
	length := len(l.input)

	// Handle negative numeric literal.
	ch := l.peek(index)
	if ch == "-" {
		value += ch
		index++
		ch = l.peek(index)
	}

	// Numbers must start with a decimal digit or a dot.
	if ch != "." && !isDecimalDigit(ch) {
		return nil
	}

	if ch != "." {
		value += l.peek(index)
		index++
		ch = l.peek(index)

		if value == "0" {
			// Base-16 numbers.
			if ch == "x" || ch == "X" {
				value += ch
				index++
				// Read hex digits.
				for index < length {
					ch = l.peek(index)
					if !isHexDigit(ch) {
						break
					}
					value += ch
					index++
				}
				if len(value) <= 2 {
					// Only "0x" found.
					return &AstNode{
						Type:        "number",
						Value:       value,
						IsMalformed: true,
						Pos:         l.Char,
					}
				}
				if index < length {
					ch = l.peek(index)
					if isIdentifierStart(ch) {
						return nil
					}
				}
				return &AstNode{
					Type:        "number",
					Value:       value,
					Base:        16,
					IsMalformed: false,
					Pos:         l.Char,
				}
			}
			// Base-8 numbers.
			if isOctalDigit(ch) {
				value += ch
				index++
				bad := false
				for index < length {
					ch = l.peek(index)
					// If a decimal digit (and not an octal digit) is encountered, mark as malformed.
					if isDecimalDigit(ch) {
						bad = true
					}
					if !isOctalDigit(ch) {
						// If not a valid punctuator, abort.
						if !isPunctuator(ch) {
							return nil
						}
						break
					}
					value += ch
					index++
				}
				if index < length {
					ch = l.peek(index)
					if isIdentifierStart(ch) {
						return nil
					}
				}
				// For simplicity, we store the value as-is.
				return &AstNode{
					Type:        "number",
					Value:       value,
					Base:        8,
					IsMalformed: bad,
				}
			}
			// Decimal numbers starting with 0 followed by a digit (e.g. "09") are considered malformed.
			if isDecimalDigit(ch) {
				value += ch
				index++
			}
		}
		// Continue reading decimal digits.
		for index < length {
			ch = l.peek(index)
			if !isDecimalDigit(ch) {
				break
			}
			value += ch
			index++
		}
	}

	// Fractional part.
	if ch == "." {
		value += ch
		index++
		for index < length {
			ch = l.peek(index)
			if !isDecimalDigit(ch) {
				break
			}
			value += ch
			index++
		}
	}

	// Exponent part.
	if ch == "e" || ch == "E" {
		value += ch
		index++
		ch = l.peek(index)
		if ch == "+" || ch == "-" {
			value += ch
			index++
		}
		ch = l.peek(index)
		if isDecimalDigit(ch) {
			value += ch
			index++
			for index < length {
				ch = l.peek(index)
				if !isDecimalDigit(ch) {
					break
				}
				value += ch
				index++
			}
		} else {
			return nil
		}
	}

	// Optionally, if the next character is not a punctuator then fail.
	if index < length {
		ch = l.peek(index)
		if !isPunctuator(ch) {
			return nil
		}
	}

	return &AstNode{
		Type:  "number",
		Value: value,
		Base:  10,
		Pos:   l.Char,
		IsMalformed: func() bool {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return true
			}
			return !isFinite(v)
		}(),
	}
}

// scanPunctuator extracts a punctuator token.
func (l *Lexer) scanPunctuator() *AstNode {
	ch := l.peek(0)
	if isPunctuator(ch) {
		return &AstNode{
			Type:  ch,
			Value: ch,
			Pos:   l.Char,
		}
	}
	return nil
}

// scanStringLiteral extracts a string literal token.
func (l *Lexer) scanStringLiteral() *AstNode {
	quote := l.peek(0)

	// String must start with a quote.
	if quote != `"` && quote != "'" {
		return nil
	}
	// Skip the opening quote.
	l.skip(1)
	var value string

	for {
		ch := l.peek(0)
		if ch == "" {
			// End of input before closing quote.
			return &AstNode{
				Type:       "string",
				Value:      value,
				IsUnclosed: true,
				Quote:      quote,
				Pos:        l.Char,
			}
		}

		if ch == quote {
			// Found closing quote.
			l.skip(1)
			break
		}

		value += ch
		l.skip(1)
	}
	return &AstNode{
		Type:       "string",
		Value:      value,
		IsUnclosed: false,
		Quote:      quote,
		Pos:        l.Char,
	}
}
