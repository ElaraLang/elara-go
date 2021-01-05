package util

import (
	"fmt"
	"reflect"
)

func NillableStringify(nillableStr *string, defaultStr string) string {
	if nillableStr == nil {
		return defaultStr
	}
	return *nillableStr
}
func Stringify(s interface{}) string {
	switch t := s.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case float64:
		return fmt.Sprintf("%g", t)
	case bool:
		return fmt.Sprintf("%t", t)
	}

	return fmt.Sprintf("%s: %s", reflect.TypeOf(s), s)
}
