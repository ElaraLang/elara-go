package interpreter

import "strings"

//In a proper implementation these would be persistent. But for now, they will do
type Collection struct {
	ElementType Type
	Elements    []*Value
}

type CollectionType struct {
	ElementType Type
}

func (t *CollectionType) Name() string {
	return t.ElementType.Name() + "[]" //Eg String[]
}

func (t *CollectionType) Accepts(other Type) bool {
	otherColl, ok := other.(*CollectionType)
	if !ok {
		return false
	}

	return t.ElementType.Accepts(otherColl.ElementType)
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
	elemStrings := make([]string, len(t.Elements))
	for i, element := range t.Elements {
		elemStrings[i] = *element.String()
	}
	return "[" + strings.Join(elemStrings, ", ") + "]"
}
