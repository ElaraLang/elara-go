package parser

import "github.com/ElaraLang/elara/lexer"

// TokenTape represents an intermediate structure between the lexer and parser
// It handles reading from lexer through channel if needed
// Channel is the channel it would listen to if isRepl is true
type TokenTape struct {
	Channel chan lexer.Token
	tokens  []lexer.Token
	index   int
	isRepl  bool
}

// tokenAt returns the token at specified index
// attempts to read tokens from channel if And only if isRepl is true and index is not on tape
func (tStream *TokenTape) tokenAt(index int) lexer.Token {
	if index > len(tStream.tokens) {
		if !tStream.isRepl {
			return lexer.CreateBlankToken(lexer.EOF)
		}
		// If in a REPL, try to read further from the channel
		required := index - len(tStream.tokens)
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

// Current returns token at current tape head
func (tStream *TokenTape) Current() lexer.Token {
	return tStream.tokenAt(tStream.index)
}

// Consume attempts to match current token with any of the provided token types and returns the same
// It fails with a parser error if none found
func (tStream *TokenTape) Consume(tokenType ...lexer.TokenType) lexer.Token {
	cur := tStream.Current()
	tStream.advance()
	notFound := true
	for _, typ := range tokenType {
		if cur.TokenType == typ {
			notFound = false
			break
		}
	}
	if notFound {
		// panic()
	}
	return cur
}

// Expect functions exactly the same as Consume but without returning consumed token
// Exists for the sake of readability
func (tStream *TokenTape) Expect(tokenType ...lexer.TokenType) {
	_ = tStream.Consume(tokenType...)
}

// FindDepthClosingIndex finds index of the closing type provided at the same depth
func (tStream *TokenTape) FindDepthClosingIndex(opening lexer.TokenType, closing lexer.TokenType) int {
	tStream.Expect(opening)
	offset := 0
	depth := 1
	for {
		switch tStream.tokenAt(tStream.index + offset).TokenType {
		case opening:
			depth++
		case closing:
			depth--
		case lexer.EOF:
			// panic
		}
		if depth == 0 {
			break
		}
		offset++
	}
	return tStream.index + offset
}
