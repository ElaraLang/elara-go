lexer grammar ElaraLexer;


NewLine : ('\r' '\n' | '\n' | '\r') {
    this.processNewLine();
} ;

Tab : [\t]+ {
    this.processTabToken();
} ;

Semicolon : ';';

fragment Space: [ \t];
Whitespace : Space+ -> skip;

InlineComment : '//' ~[\r\n]* -> skip;
MultiComment : '/*' .+? '*/' -> skip;

fragment NonZeroDigit : '1'..'9';
fragment Digit: '0' | NonZeroDigit;
fragment HexDigit: [0-9A-F];
fragment ScientificNotation: 'E' [+-];

// Literals
fragment AbsoluteIntegerLiteral: NonZeroDigit (Digit+ (ScientificNotation Digit+))?;
IntegerLiteral : '-'? AbsoluteIntegerLiteral;
FloatLiteral: IntegerLiteral '.' AbsoluteIntegerLiteral ;

CharLiteral : '\'' (. | '\\' .) '\'';
StringLiteral : '"' (~'"' | '\\"')+ '"';

// Keywords
Let : 'let';
Def : 'def';
Mut: 'mut';
Type: 'type';
Class: 'class';

// Symbols
Comma: ',';
LParen: '(';
RParen: ')';
LSquareParen: '[';
RSquareParen: ']';
LBrace: '{';
RBrace: '}';
Colon: ':';
Dot: '.';
PureArrow: '->';
ImpureArrow: '=>';
Equals : '=';
Bar : '|';

// Identifiers
TypeIdentifier: [A-Z][a-zA-Z_0-9]*;
VarIdentifier: [a-z][a-zA-Z_0-9]*;
OperatorIdentifier: ('!' | '#' | '$' | '%' | '+' | '-' | '/' | '*' | '.' | '<' | '>' | '=' | '?' | '@' | '~' | '\\' | '^' | '|')+;

