package byte

import (
	"testing"
	"bytes"
	"fmt"
)

func TestInsertThatRequiresReSlice(testCtx *testing.T) {
	// given
	var expected = []byte("this is to test the amazingness of the insert function")
	slice := make([]byte, 32*1024)
	sliceData := []byte("this is to test the of the insert function")
	copy(slice, sliceData)

	searchString := "of the insert"
	insertLocation := bytes.Index(slice[0:len(sliceData)], []byte(searchString))
	insertString := []byte("amazingness ")

	// when
	actual := Insert(slice[0:len(sliceData)], insertLocation, insertString)

	// then
	if !bytes.Equal(expected, actual) {
		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expected, actual))
	}
}

func TestInsertThatRequiresNewSlice(testCtx *testing.T) {
	// given
	var expected = []byte("this is to test the amazingness of the insert function")
	slice := []byte("this is to test the of the insert function")

	searchString := "of the insert"
	insertLocation := bytes.Index(slice, []byte(searchString))
	insertString := []byte("amazingness ")

	// when
	actual := Insert(slice, insertLocation, insertString)

	// then
	if !bytes.Equal(expected, actual) {
		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expected, actual))
	}
}

func TestInsertDataLongerThenOriginal(testCtx *testing.T) {
	// given
	var expected = []byte("this is to test the amazingness and greatness of the this very very long text to test inserting somthing longer then the original sentence")
	slice := []byte("this is to test the sentence")

	searchString := "sentence"
	insertLocation := bytes.Index(slice, []byte(searchString))
	insertString := []byte("amazingness and greatness of the this very very long text to test inserting somthing longer then the original ")

	// when
	actual := Insert(slice, insertLocation, insertString)

	// then
	if !bytes.Equal(expected, actual) {
		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expected, actual))
	}
}

func TestInsertForEmptyData(testCtx *testing.T) {
	// given
	var expected = []byte("this is to test the insert function")
	slice := []byte("this is to test the insert function")

	searchString := "insert"
	insertLocation := bytes.Index(slice, []byte(searchString))
	insertString := []byte("")

	// when
	actual := Insert(slice, insertLocation, insertString)

	// then
	if !bytes.Equal(expected, actual) {
		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expected, actual))
	}
}
//
//func TestInsertForIndexTooLarge(testCtx *testing.T) {
//	// given
//	var expected = []byte("this is to test the insert function")
//	data := make([]byte, 32*1024)
//	data = []byte("this is to test the insert function")
//
//	insertLocation := len(data) + 5
//	insertString := []byte("")
//
//	// when
//	actual := Insert(data, insertLocation, insertString)
//
//	// then
//	if !bytes.Equal(expected, actual) {
//		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expected, actual))
//	}
//}

