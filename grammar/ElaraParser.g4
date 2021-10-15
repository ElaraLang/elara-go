parser grammar ElaraParser;

options { tokenVocab=ElaraLexer; }

defClause : Def VarIdentifier Colon type;
letClause : Let VarIdentifier Equals expression;
variable : defClause? letClause;

unit : LParen RParen;

// Types
type : unit #UnitType
    | VarIdentifier #GenericType
    | typeName typeName+ #RealizedGenericType
    | type PureArrow type # PureFunctionType
    | type ImpureArrow type #ImpureFunctionType
    | LParen type RParen #ParenthesizedType
    | LParen (type (Comma type)+) RParen #TupleType
    | LSquareParen type RSquareParen # ListType
    | TypeIdentifier #SimpleType
    | Mut type #MutType;

typeName :
      TypeIdentifier #NormalTypeName
    | VarIdentifier #GenericTypeName;

typeAlias : type;
typeConstructor :
    TypeIdentifier type* #NormalTypeConstructor
    | TypeIdentifier recordType #RecordTypeConstructor;

sumType : typeConstructor (Bar typeConstructor)*;
recordTypeField : VarIdentifier Colon type;
recordType : LBrace recordTypeField (Comma recordTypeField)* RBrace;

typeDeclaration : Type TypeIdentifier VarIdentifier* Equals typeDeclarationValue;
typeDeclarationValue : typeAlias | sumType | recordType;

// Expressions

expression :
    unit #UnitExpression
    | IntegerLiteral #IntExpression
    | FloatLiteral #FloatExpression
    | CharLiteral #CharExpression
    | StringLiteral #StringExpression
    | LSquareParen (expression (Comma expression)*)? RSquareParen #ListExpression
    | LParen (expression (Comma expression)+) RParen #TupleExpression
;