package util

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"hash/fnv"
	"reflect"
	"strconv"
)

func NillableStringify(nillableStr *string, defaultStr string) string {
	if nillableStr == nil {
		return defaultStr
	}
	return *nillableStr
}

func Stringify(s interface{}) string {
	switch t := s.(type) {
	case Stringable:
		return t.String()
	case string:
		return t
	case rune:
		return string(t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	}

	return fmt.Sprintf("%s: %s", reflect.TypeOf(s), s)
}

type Stringable interface {
	String() string
}

func Hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

func JoinToString(slice interface{}, separator string) string {
	res := ""

	switch slice := slice.(type) {
	case []string:
		l := len(slice)
		for i, v := range slice {
			res += v
			if i < l {
				res += separator
			}
		}
	case []ast.Entry:
		l := len(slice)
		for i, v := range slice {
			res += "(" + v.Key.ToString() + ") :" + "(" + v.Value.ToString() + ")"
			if i < l {
				res += separator
			}
		}

	case []ast.Parameter:
		l := len(slice)
		for i, v := range slice {
			res += v.ToString()
			if i < l {
				res += separator
			}
		}
	case []ast.Node:
		l := len(slice)
		for i, v := range slice {
			res += v.ToString()
			if i < l {
				res += separator
			}
		}
	default:
		res = "Unknown"
	}
	return res
}
