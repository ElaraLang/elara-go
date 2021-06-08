package util

import (
	"fmt"
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
