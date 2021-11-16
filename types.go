package AES

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	Int8    = "int8"
	Uint8   = "uint8"
	Byte    = "[]uint8"
	Int16   = "int16"
	Uint16  = "uint16"
	Int32   = "int32"
	Rune    = "rune"
	Uint32  = "uint32"
	Int64   = "int64"
	Uint64  = "uint64"
	Int     = "int"
	Uint    = "uint"
	Uintptr = "uintptr"
	Float32 = "float32"
	Float64 = "float64"
	String  = "string"
	Bool    = "bool"
)

func Convert(s string, t reflect.Value) (interface{}, error) {
	switch t.Type().String() {
	case Int8:
		var i int8

		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}

		i = int8(v)

		return i, nil

	case Uint8:
		var u uint8

		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}

		u = uint8(v)

		return u, nil

	case Byte:
		var b []byte

		for _, v := range s {
			val := string(v)
			if val == "[" || val == "]" || val == " " {
				continue
			}
			ival, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("could not convert %s to string while structuring byte array", string(v))
			}

			b = append(b, uint8(ival))
		}

		return b, nil

	case Int16:
		var i int16

		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}

		i = int16(v)

		return i, nil

	case Uint16:
		var u uint16

		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}

		u = uint16(v)

		return u, nil

	case Int32:
		var i int32

		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}

		i = int32(v)

		return i, nil

	case Uint32:
		var u uint32

		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}

		u = uint32(v)

		return u, nil

	case Int64:
		var i int64

		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}

		i = int64(v)

		return i, nil

	case Uint64:
		var u uint64

		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}

		u = uint64(v)

		return u, nil

	case Int:
		var i int

		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		i = int(v)

		return i, nil

	case Uint:
		var u uint

		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		u = uint(v)

		return u, nil

	case Uintptr:
		var u uintptr

		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		u = uintptr(v)

		return u, nil

	case Float32:
		var f float32

		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}

		f = float32(v)

		return f, nil

	case Float64:
		var f float64

		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}

		return f, nil

	case String:
		return s, nil

	case Bool:
		var b bool

		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}

		return b, nil

	default:
		return nil, fmt.Errorf("%s is not a supported file type", t.Type().String())
	}
}
