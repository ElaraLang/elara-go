package parser

import "github.com/ElaraLang/elara/lexer"

// TokenTape represents an intermediate structure between the lexer and Parser
// It handles reading from lexer through channel if needed
// Channel is the channel it would listen to if isRepl is true
type TokenTape struct {
	Channel chan lexer.Token
	tokens  []lexer.Token
	index   int
}

// NewTokenTape creates a TokenTape that manages tokens from input channel
func NewTokenTape(channel chan lexer.Token) TokenTape {
	return TokenTape{
		Channel: channel,
		tokens:  []lexer.Token{},
		index:   0,
	}
}

// IsClosed checks where the token will continue to read
// If true, use Unwind to open the tape for writing
func (tStream *TokenTape) IsClosed() bool {
	if len(tStream.tokens) < 1 {
		return false
	}
	return tStream.tokens[len(tStream.tokens)-1].TokenType == lexer.EOF
}

// Unwind cleans cached tokens and resets tape head
func (tStream *TokenTape) Unwind() {
	tStream.tokens = []lexer.Token{}
	tStream.index = 0
}

// skipLineBreaks moves head to the first non-NEWLINE token with index >= current index
func (tStream *TokenTape) skipLineBreaks() {
	for tStream.Current().TokenType == lexer.NEWLINE {
		tStream.advance()
	}
}

// tokenAt returns the token at specified index
// attempts to read tokens from channel if And only if isRepl is true and index is not on tape
func (tStream *TokenTape) tokenAt(index int) lexer.Token {
	if index >= len(tStream.tokens) {
		if len(tStream.tokens) > 1 {
			lastToken := tStream.tokenAt(len(tStream.tokens) - 1)
			if lastToken.TokenType == lexer.EOF {
				return lastToken
			}
		}
		required := index - len(tStream.tokens) + 1
		tStream.readFromChannel(required)
	}
	return tStream.tokens[index]
}

// readFromChannel attempts to read specified amount of tokens from the tape's Channel
func (tStream *TokenTape) readFromChannel(amount int) {
	for amount > 0 {
		amount--
		tok := <-tStream.Channel
		tStream.Append(tok)
	}
}

// moveHead moves the head of the tape by a specific amount.
// Exists for the sake of readability
func (tStream *TokenTape) moveHead(amount int) {
	tStream.index += amount
}

// advance moves the tape head forward by 1
func (tStream *TokenTape) advance() {
	tStream.moveHead(1)
}

// Append appends provided tokens to the end of the token stream
func (tStream *TokenTape) Append(inputTokens ...lexer.Token) {
	tStream.tokens = append(tStream.tokens, inputTokens...)
}

// Peek returns the token at an offset of specified amount from current index
func (tStream *TokenTape) Peek(amount int) lexer.Token {
	return tStream.tokenAt(tStream.index + amount)
}

// ValidateHead validates if any of the provided token is at current head
func (tStream *TokenTape) ValidateHead(tokenType ...lexer.TokenType) bool {
	return tStream.ValidationPeek(0, tokenType...)
}

// ValidationPeek peeks by given amount and validates if the provided token is at that index
func (tStream *TokenTape) ValidationPeek(amount int, tokenType ...lexer.TokenType) bool {
	actual := tStream.Peek(amount).TokenType
	for _, v := range tokenType {
		if actual == v {
			return true
		}
	}
	return false
}

// Current returns token at current tape head
func (tStream *TokenTape) Current() lexer.Token {
	return tStream.tokenAt(tStream.index)
}

// Validate current token and advance if successful
func (tStream *TokenTape) Match(tokenType ...lexer.TokenType) bool {
	found := tStream.ValidationPeek(0, tokenType...)
	if found {
		tStream.advance()
	}
	return found
}

// ConsumeAny reads the current token independent of its type
func (tStream *TokenTape) ConsumeAny() lexer.Token {
	cur := tStream.Current()
	tStream.advance()
	return cur
}

// Consume attempts to match current token with any of the provided token types and returns the same
// It fails with a Parser error if none found
func (tStream *TokenTape) Consume(tokenType ...lexer.TokenType) lexer.Token {
	cur := tStream.Current()
	found := false
	for _, typ := range tokenType {
		if cur.TokenType == typ {
			found = true
			break
		}
	}
	if !found {
		panic(NewParseError(cur, "Invalid token consumed by tape!"))
	}
	tStream.advance()
	return cur
}

func (tStream *TokenTape) MatchUnorderedSequence(tokenType ...lexer.TokenType) map[lexer.TokenType]bool {
	result := map[lexer.TokenType]bool{}
	for _, v := range tokenType {
		result[v] = false
	}
	for tStream.ValidationPeek(0, tokenType...) {
		result[tStream.Current().TokenType] = true
		tStream.advance()
	}
	return result
}

// Expect functions exactly the same as Consume but without returning consumed token
// Exists for the sake of readability
func (tStream *TokenTape) Expect(tokenType ...lexer.TokenType) {
	_ = tStream.Consume(tokenType...)
}

// FindDepthClosingIndex finds index of the closing type provided at the same depth
func (tStream *TokenTape) FindDepthClosingIndex(opening lexer.TokenType, closing lexer.TokenType) int {
	if !tStream.ValidateHead(opening) {
		return tStream.index
	}
	offset := 0
	depth := 1
	for {
		switch tStream.tokenAt(tStream.index + offset).TokenType {
		case opening:
			depth++
		case closing:
			depth--
		case lexer.EOF:
			panic(NewParseError(tStream.tokenAt(tStream.index+offset), "Token Stream ended unexpectedly!"))
		}
		if depth == 0 {
			break
		}
		offset++
	}
	return tStream.index + offset
}
