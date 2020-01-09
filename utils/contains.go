package utils

// ContainsUint32 return true if sub is present in arr
func ContainsUint32(arr []uint32, sub uint32) bool {
	for _, a := range arr {
		if a == sub {
			return true
		}
	}

	return false
}

// ContainsByte return true if sub is present in arr
func ContainsByte(arr []byte, sub byte) bool {
	for _, a := range arr {
		if a == sub {
			return true
		}
	}

	return false
}
