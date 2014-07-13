package byte

func Insert(slice []byte, index int, value []byte) []byte {
	if index < cap(slice) {
		insertLength := len(value)
		if (len(slice)+insertLength) > cap(slice) {
			slice = expandSlice(slice, len(slice)+insertLength)
		} else {
			slice = slice[0:len(slice)+insertLength]
		}

		copy(slice[index+insertLength:], slice[index:])
		for i := 0; i < insertLength; i++ {
			slice[index+i] = value[i]
		}
	}
	return slice
}

func expandSlice(slice []byte, newLength int) []byte {
	newSlice := make([]byte, newLength)
	copy(newSlice, slice)
	return newSlice
}



