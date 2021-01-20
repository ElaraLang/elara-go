package parser

const (
	_ int = iota
	LOWEST
	EQUALS     // ==
	COMPARISON // > or <
	SUM        // +
	PRODUCT    // *
	PREFIX     // -X or !X
	CALL       // myFunction(X)
)
