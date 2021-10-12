parser grammar ElaraParser;

options { tokenVocab=ElaraLexer; }

defClause : Def VarIdentifier Colon type;
letClause : Let VarIdentifier Equals expression;
variable : (defClause NewLine)? letClause;

unit : LParen RParen;


type : unit #UnitType
    | VarIdentifier #GenericType
    | type PureArrow type # PureFunctionType
    | type ImpureArrow type #ImpureFunctionType
    | LParen type RParen #ParenthesizedType
    | TypeIdentifier #SimpleType;

typeAlias : Type TypeIdentifier Equals type;
typeConstructor : TypeIdentifier type*;

sumType : typeConstructor (Bar typeConstructor)*;

typeDeclaration : typeAlias | sumType;

expression : unit;
