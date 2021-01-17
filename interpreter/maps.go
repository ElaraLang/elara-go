package interpreter

import (
	"fmt"
	"strings"
)

type MapType struct {
	KeyType   Type
	ValueType Type
}

type Map struct {
	MapType  *MapType
	Elements []*Entry
}

type Entry struct {
	Key   *Value
	Value *Value
}

func (m *Map) Get(ctx *Context, key *Value) *Value {
	for _, element := range m.Elements {
		if element.Key.Equals(ctx, key) {
			return element.Value
		}
	}
	return nil
}
func MapOf(elements []*Entry) *Map {
	mapType := &MapType{
		KeyType:   AnyType,
		ValueType: AnyType, //tfw type inference
	}
	return &Map{
		MapType:  mapType,
		Elements: elements,
	}
}

func (t *MapType) Name() string {
	return fmt.Sprintf("{ %s : %s }", t.KeyType.Name(), t.ValueType.Name())
}

func (t *MapType) Accepts(otherType Type, ctx *Context) bool {
	asMap, isMap := otherType.(*MapType)
	if !isMap {
		return false
	}
	return t.KeyType.Accepts(asMap.KeyType, ctx) && t.ValueType.Accepts(asMap.ValueType, ctx)
}

func (m *Map) String() string {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for i, elem := range m.Elements {
		builder.WriteRune('\n')
		builder.WriteString(elem.Key.String())
		builder.WriteString(": ")
		builder.WriteString(elem.Value.String())
		if i != len(m.Elements)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune('\n')
	builder.WriteRune('}')
	return builder.String()
}
