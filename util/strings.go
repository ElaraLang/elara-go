package util

import (
	"fmt"
	"reflect"
)

func Stringify(s interface{}) string {
	switch t := s.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case bool:
		return fmt.Sprintf("%t", t)
	}

	return fmt.Sprintf("%s: %s", reflect.TypeOf(s), s)
}
