package byte

//import "fmt"

func Insert(slice []byte, index int, value []byte) []byte {
	insertLength := len(value)
	if (len(slice)+insertLength) > cap(slice) {
		slice = expandSlice(slice, len(slice)+insertLength)
	} else {
		slice = slice[0:len(slice)+insertLength]
	}

//	fmt.Printf("Before insert slice: \n%s\n", slice)
	copy(slice[index+insertLength:], slice[index:])
	for i := 0; i < insertLength; i++ {
		slice[index+i] = value[i]
	}
//	fmt.Printf("After insert slice: \n%s\n", slice)

	return slice
}

func expandSlice(slice []byte, newLength int) []byte {
	newSlice := make([]byte, newLength)
	copy(newSlice, slice)
	return newSlice
}



