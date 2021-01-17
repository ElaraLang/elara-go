package interpreter

import "strings"

//In a proper implementation these would be persistent. But for now, they will do
type Collection struct {
	ElementType Type
	Elements    []*Value

	cachedAsString *string
}

type CollectionType struct {
	ElementType Type
}

func (t *CollectionType) Name() string {
	return "[" + t.ElementType.Name() + "]" //Eg [String]
}

func (t *CollectionType) Accepts(otherType Type, ctx *Context) bool {
	otherColl, ok := otherType.(*CollectionType)
	if !ok {
		return false
	}

	return t.ElementType.Accepts(otherColl.ElementType, ctx)
}

func NewCollectionType(collection *Collection) Type {
	return &CollectionType{
		collection.ElementType,
	}
}
func NewCollectionTypeOf(elemType Type) Type {
	return &CollectionType{
		elemType,
	}
}

func (t *Collection) String() string {
	if t.ElementType == CharType {
		return t.elemsAsString()
	}
	elemStrings := make([]string, len(t.Elements))
	for i, element := range t.Elements {
		elemStrings[i] = element.String()
	}
	return "[" + strings.Join(elemStrings, ", ") + "]"
}

func (t *Collection) elemsAsString() string {
	if t.cachedAsString != nil {
		return *t.cachedAsString
	}

	if t.ElementType != CharType {
		panic("Cannot convert collection to string")
	}
	builder := strings.Builder{}
	for _, elem := range t.Elements {
		builder.WriteRune(elem.Value.(rune))
	}

	asString := builder.String()
	t.cachedAsString = &asString
	return asString
}
