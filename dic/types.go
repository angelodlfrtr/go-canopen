package dic

import (
	"github.com/thoas/go-funk"
)

const (
	Var byte = 0x07
	Arr byte = 0x08
	Rec byte = 0x09
)

const (
	Boolean    byte = 0x1
	Integer8   byte = 0x2
	Integer16  byte = 0x3
	Integer32  byte = 0x4
	Integer64  byte = 0x15
	Unsigned8  byte = 0x5
	Unsigned16 byte = 0x6
	Unsigned32 byte = 0x7
	Unsigned64 byte = 0x1b

	Real32 byte = 0x8
	Real64 byte = 0x11

	VisibleString byte = 0x9
	OctetString   byte = 0xa
	UnicodeString byte = 0xb
	Domain        byte = 0xf
)

func IsSignedType(t byte) bool {
	return funk.Contains([]byte{
		Integer8,
		Integer16,
		Integer32,
		Integer64,
	}, t)
}

func IsUnsignedType(t byte) bool {
	return funk.Contains([]byte{
		Unsigned8,
		Unsigned16,
		Unsigned32,
		Unsigned64,
	}, t)
}

func IsIntegerType(t byte) bool {
	return funk.Contains([]byte{
		Unsigned8,
		Unsigned16,
		Unsigned32,
		Unsigned64,
		Integer8,
		Integer16,
		Integer32,
		Integer64,
	}, t)
}

func IsFloatType(t byte) bool {
	return funk.Contains([]byte{
		Real32,
		Real64,
	}, t)
}

func IsNumberType(t byte) bool {
	return funk.Contains([]byte{
		Unsigned8,
		Unsigned16,
		Unsigned32,
		Unsigned64,
		Integer8,
		Integer16,
		Integer32,
		Integer64,
		Real32,
		Real64,
	}, t)
}

func IsDataType(t byte) bool {
	return funk.Contains([]byte{
		VisibleString,
		OctetString,
		UnicodeString,
		Domain,
	}, t)
}
