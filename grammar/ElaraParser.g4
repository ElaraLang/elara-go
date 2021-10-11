parser grammar ElaraParser;

options { tokenVocab=ElaraLexer; }

defClause : Def VarIdentifier Colon type;
letClause : Let VarIdentifier Equals expression;
variable : (defClause NewLine)? letClause;

unitValue : LParen RParen;
unitType : LParen RParen;

type : unitType
    | VarIdentifier // Generic types
    | TypeIdentifier;

expression : unitValue;
