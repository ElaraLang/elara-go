parser grammar ElaraParser;

options { tokenVocab=ElaraLexer; }

defClause : Def VarIdentifier Colon type;
letClause : Let VarIdentifier Equals expression;
variable : (defClause NewLine)? letClause;

unit : LParen RParen;


type : unit #UnitType
    | VarIdentifier #GenericType
    | typeName typeName+ #RealizedGenericType
    | type PureArrow type # PureFunctionType
    | type ImpureArrow type #ImpureFunctionType
    | LParen type RParen #ParenthesizedType
    | TypeIdentifier #SimpleType
    | Mut type #MutType;

typeName :
      TypeIdentifier #NormalTypeName
    | VarIdentifier #GenericTypeName;


typeAlias : type;
typeConstructor :
    TypeIdentifier type* #NormalTypeCosntructor
    | TypeIdentifier recordType #RecordTypeConstructor;

sumType : typeConstructor (Bar typeConstructor)*;
recordTypeField : VarIdentifier Colon type;
recordType : LBrace recordTypeField (Comma recordTypeField)* RBrace;

typeDeclaration : Type TypeIdentifier VarIdentifier* Equals typeDeclarationValue;
typeDeclarationValue : typeAlias | sumType | recordType;

expression : unit;
