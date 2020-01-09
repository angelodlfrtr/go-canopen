package utils

import (
	"testing"
)

func TestContainsUint32(t *testing.T) {
	arr := []uint32{0x01, 0x02, 0x03}

	if !ContainsUint32(arr, uint32(0x03)) {
		t.Fatalf("ContainsUint32 with %v / %v should return true", arr, 0x03)
	}
}

func TestContainsByte(t *testing.T) {
	arr := []byte{0x01, 0x02, 0x03}

	if !ContainsByte(arr, byte(0x03)) {
		t.Fatalf("ContainsByte with %v / %v should return true", arr, 0x03)
	}
}
