package assertion

import (
	"testing"
	"bytes"
	"reflect"
	"fmt"
)

type visit struct {
	a1  uintptr
	a2  uintptr
	typ reflect.Type
}

const equal_comparison_failure_message = "\n%s - Failed: - %s - expected: [%v] not equal to actual: [%v]\n"

func AssertArrayEquals(testCtx *testing.T, expected []byte, actual []byte) {
	for key := range expected {
		if expected[key] != actual[key] {
			testCtx.Fatalf("\nFailed index [%d] value [%s] %v not equal to [%s] %v\n\n\n", key, expected[key:key+1], expected[key:key+1], actual[key:key+1], actual[key:key+1])
		} else {
			testCtx.Logf("Index [%d] value [%s] %v is equal to [%s] %v", key, expected[key:key+1], expected[key:key+1], actual[key:key+1], actual[key:key+1])
		}
	}
}

func AssertDeepEqual(message string, testCtx *testing.T, expected, actual interface{}) {
	failure_message := equal_comparison_failure_message + fmt.Sprintf("\t expected:\n\t[%#v]\n\t[%s]\n\t actual:\n\t[%#v]\n\t[%s]\n\n\n", expected, expected, actual, actual)
	if expected != nil && actual != nil {
		v1 := reflect.ValueOf(expected)
		v2 := reflect.ValueOf(actual)
		if v1.Type() != v2.Type() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "type not equal", v1.Type(), v2.Type()))
		}
		deepValueEqual(failure_message, message, testCtx, v1, v2, make(map[visit]bool))
	} else if expected != actual {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "one value is nil", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual)))
	}
}

func deepValueEqual(failure_message, message string, testCtx *testing.T, expected, actual reflect.Value, visited map[visit]bool) {

	if !expected.IsValid() || !actual.IsValid() {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "not valid value", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual)))
	}
	if expected.Type() != actual.Type() {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "type not equal", expected.Type(), actual.Type()))
	}

	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			return true
		}
		return false
	}

	if expected.CanAddr() && actual.CanAddr() && hard(expected.Kind()) {
		addr1 := expected.UnsafeAddr()
		addr2 := actual.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are identical ...
		if addr1 == addr2 {
			return
		}

		// ... or already seen
		typ := expected.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return
		}

		// Remember for later.
		visited[v] = true
	}

	switch expected.Kind() {
	case reflect.Array:
		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "length not equal", expected.Len(), actual.Len()))
		}
		for i := 0; i < expected.Len(); i++ {
			deepValueEqual(failure_message, message, testCtx, expected.Index(i), actual.Index(i), visited)
		}
	case reflect.Slice:
		if expected.IsNil() != actual.IsNil() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "one slice is nil and the other is not", fmt.Sprintf("is nil: %t", expected.IsNil()), fmt.Sprintf("is nil: %t", actual.IsNil())))
		}
		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "length not equal", expected.Len(), actual.Len()))
		}
		if expected.Pointer() != actual.Pointer() {
			for i := 0; i < expected.Len(); i++ {
				deepValueEqual(failure_message, message, testCtx, expected.Index(i), actual.Index(i), visited)
			}
		}
	case reflect.Interface:
		if !expected.IsNil() && !actual.IsNil() {
			deepValueEqual(failure_message, message, testCtx, expected.Elem(), actual.Elem(), visited)
		} else if !(expected.IsNil() && actual.IsNil()) {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "objects not equal", expected, actual))
		}
	case reflect.Ptr:
		var nilPointer uintptr
		if expected.Pointer() != nilPointer && actual.Pointer() != nilPointer {
			deepValueEqual(failure_message, message, testCtx, expected.Elem(), actual.Elem(), visited)
		} else if expected.Pointer() != actual.Pointer() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "one pointer is nil and the other is not", fmt.Sprintf("%#v", expected.Pointer()), fmt.Sprintf("%#v", actual.Pointer())))
		}
	case reflect.Struct:
		for i, n := 0, expected.NumField(); i < n; i++ {
			deepValueEqual(failure_message, message, testCtx, expected.Field(i), actual.Field(i), visited)
		}
	case reflect.Map:
		if expected.IsNil() != actual.IsNil() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "one map is nil and the other is not", fmt.Sprintf("is nil: %t", expected.IsNil()), fmt.Sprintf("is nil: %t", actual.IsNil())))
		}
		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "keys not equal", expected.MapKeys(), actual.MapKeys()))
		}
		if expected.Pointer() != actual.Pointer() {
			for _, k := range expected.MapKeys() {
				deepValueEqual(failure_message, message, testCtx, expected.MapIndex(k), actual.MapIndex(k), visited)
			}
		}
	case reflect.Func:
		if !(expected.IsNil() && actual.IsNil()) {
			// Can't do better than this:
			testCtx.Fatal(fmt.Errorf(failure_message, message, "functions not equal", expected, actual))
		}
	case reflect.Bool:
		if expected.Bool() != actual.Bool() {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "booleans not equal", expected.Bool(), actual.Bool()))
		}
	case reflect.Int:
		compareInt(failure_message, message, testCtx, expected, actual)
	case reflect.Int8:
		compareInt(failure_message, message, testCtx, expected, actual)
	case reflect.Int16:
		compareInt(failure_message, message, testCtx, expected, actual)
	case reflect.Int32:
		compareInt(failure_message, message, testCtx, expected, actual)
	case reflect.Int64:
		compareInt(failure_message, message, testCtx, expected, actual)
	case reflect.Uint:
		compareUInt(failure_message, message, testCtx, expected, actual)
	case reflect.Uint8:
		compareUInt(failure_message, message, testCtx, expected, actual)
	case reflect.Uint16:
		compareUInt(failure_message, message, testCtx, expected, actual)
	case reflect.Uint32:
		compareUInt(failure_message, message, testCtx, expected, actual)
	case reflect.Uint64:
		compareUInt(failure_message, message, testCtx, expected, actual)
	case reflect.Float32:
		compareFloat(failure_message, message, testCtx, expected, actual)
	case reflect.Float64:
		compareFloat(failure_message, message, testCtx, expected, actual)
	default:
		// Normal equality suffices
		if !bytes.Equal([]byte(fmt.Sprintf("%v", expected)), []byte(fmt.Sprintf("%v", actual))) {
			testCtx.Fatal(fmt.Errorf(failure_message, message, "values not equal", fmt.Sprintf("%v", expected), fmt.Sprintf("%v", actual)))
		}
	}
}

func compareUInt(failure_message, message string, testCtx *testing.T, expected, actual reflect.Value) {
	if expected.Uint() != actual.Uint() {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "integers not equal", expected.Uint(), actual.Uint()))
	}
}

func compareInt(failure_message, message string, testCtx *testing.T, expected, actual reflect.Value) {
	if expected.Int() != actual.Int() {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "integers not equal", expected.Int(), actual.Int()))
	}
}

func compareFloat(failure_message, message string, testCtx *testing.T, expected, actual reflect.Value) {
	if expected.Float() != actual.Float() {
		testCtx.Fatal(fmt.Errorf(failure_message, message, "floating point number not equal", expected.Float(), actual.Float()))
	}
}
