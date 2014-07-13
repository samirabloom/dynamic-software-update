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

const equal_comparison_failure_message = "\n%s - Failed: %s [%v] not equal to [%v]\n\n\n"

func AssertArrayEquals(testCtx *testing.T, expected []byte, actual []byte) {
	for key := range expected {
		if expected[key] != actual[key] {
			testCtx.Fatalf("\nFailed index [%d] value [%s] %v not equal to [%s] %v\n\n\n", key, expected[key:key+1], expected[key:key+1], actual[key:key+1], actual[key:key+1])
		} else {
			testCtx.Logf("Index [%d] value [%s] %v is equal to [%s] %v", key, expected[key:key+1], expected[key:key+1], actual[key:key+1], actual[key:key+1])
		}
	}
}

func AssertDeepEqual(message string, testCtx *testing.T, expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}


	v1 := reflect.ValueOf(expected)
	v2 := reflect.ValueOf(actual)
	if v1.Type() != v2.Type() {
		testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "type not equal", v1.Type(), v2.Type()))
		return false
	}
	return deepValueEqual(message, testCtx, v1, v2, make(map[visit]bool), 0)
}

func deepValueEqual(message string, testCtx *testing.T, expected, actual reflect.Value, visited map[visit]bool, depth int) bool {

	if !expected.IsValid() || !actual.IsValid() {
		testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "not valid value", expected, actual))
		return expected.IsValid() == actual.IsValid()
	}
	if expected.Type() != actual.Type() {
		testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "type not equal", expected.Type(), actual.Type()))
		return false
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging
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
			return true
		}

		// ... or already seen
		typ := expected.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}


	switch expected.Kind() {
	case reflect.Array:

		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "length not equal", expected.Len(), actual.Len()))
			return false
		}
		for i := 0; i < expected.Len(); i++ {
			if !deepValueEqual(message, testCtx, expected.Index(i), actual.Index(i), visited, depth+1) {
				testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "arrays not equal", expected, actual))
				return false
			}
		}
		return true
	case reflect.Slice:

		if expected.IsNil() != actual.IsNil() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "one object is nil and the other is not", expected.IsNil(), actual.IsNil()))
			return false
		}
		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "length not equal", expected.Len(), actual.Len()))
			return false
		}
		if expected.Pointer() == actual.Pointer() {
			return true
		}
		for i := 0; i < expected.Len(); i++ {
			if !deepValueEqual(message, testCtx, expected.Index(i), actual.Index(i), visited, depth+1) {
				testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "objects not equal", expected, actual))
				return false
			}
		}
		return true
	case reflect.Interface:

		if expected.IsNil() || actual.IsNil() {
			return expected.IsNil() == actual.IsNil()
		}
		return deepValueEqual(message, testCtx, expected.Elem(), actual.Elem(), visited, depth+1)
	case reflect.Ptr:

		return deepValueEqual(message, testCtx, expected.Elem(), actual.Elem(), visited, depth+1)
	case reflect.Struct:

		for i, n := 0, expected.NumField(); i < n; i++ {
			if !deepValueEqual(message, testCtx, expected.Field(i), actual.Field(i), visited, depth+1) {
				testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "structs not equal", expected, actual))
				return false
			}
		}
		return true
	case reflect.Map:

		if expected.IsNil() != actual.IsNil() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "one object is nil and the other is not", expected.IsNil(), actual.IsNil()))
			return false
		}
		if expected.Len() != actual.Len() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "keys not equal", expected.MapKeys(), actual.MapKeys()))
			return false
		}
		if expected.Pointer() == actual.Pointer() {
			return true
		}
	for _, k := range expected.MapKeys() {
		if !deepValueEqual(message, testCtx, expected.MapIndex(k), actual.MapIndex(k), visited, depth+1) {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "maps not equal", expected, actual))
			return false
		}
	}
		return true
	case reflect.Func:

		if expected.IsNil() && actual.IsNil() {
			return true
		}
		// Can't do better than this:
		testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "functions not equal", expected, actual))
		return false
	case reflect.Bool:
		if expected.Bool() != actual.Bool() {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "booleans not equal", expected.Bool(), actual.Bool()))
			return false;
		} else {
			return true;
		}
	case reflect.Int:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Int8:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Int16:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Int32:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Int64:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Uint:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Uint16:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Uint32:
		return compareInt(message, testCtx, expected, actual)
	case reflect.Uint64:
		return compareInt(message, testCtx, expected, actual)
	default:
		// Normal equality suffices
		if !bytes.Equal([]byte(fmt.Sprintf("%v", expected)), []byte(fmt.Sprintf("%v", actual))) {
			testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "values not equal", fmt.Sprintf("%v", expected), fmt.Sprintf("%v", actual)))
			return false;
		} else {
			return true;
		}
	}
}

func compareInt(message string, testCtx *testing.T, expected, actual reflect.Value) bool {
	if expected.Int() != actual.Int() {
		testCtx.Fatal(fmt.Errorf(equal_comparison_failure_message, message, "integers not equal", expected.Int(), actual.Int()))
		return false;
	} else {
		return true;
	}
}
