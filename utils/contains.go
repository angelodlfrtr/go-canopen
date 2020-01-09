package utils

func ContainsUint32(arr []uint32, sub uint32) bool {
	for _, a := range arr {
		if a == sub {
			return true
		}
	}

	return false
}
