package encryption

import (
	"fmt"
	"reflect"
	"strings"
)

//Define all outlier types
const (
	StringArr = "[]string"
)

func Encode(value reflect.Value) string {
	switch value.Type().String() {
	case StringArr:
		return strings.Join(value.Interface().([]string), "°")
	default:
		return fmt.Sprint(value)
	}
}
