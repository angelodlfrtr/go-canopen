package dic

import (
	"github.com/thoas/go-funk"
)

const (
	VAR    byte = 0x07
	ARR    byte = 0x08
	RECORD byte = 0x09
)

const (
	BOOLEAN    byte = 0x1
	INTEGER8   byte = 0x2
	INTEGER16  byte = 0x3
	INTEGER32  byte = 0x4
	INTEGER64  byte = 0x15
	UNSIGNED8  byte = 0x5
	UNSIGNED16 byte = 0x6
	UNSIGNED32 byte = 0x7
	UNSIGNED64 byte = 0x1B

	REAL32 byte = 0x8
	REAL64 byte = 0x11

	VISIBLE_STRING byte = 0x9
	OCTET_STRING   byte = 0xA
	UNICODE_STRING byte = 0xB
	DOMAIN         byte = 0xF
)

func IsSignedType(t byte) bool {
	return funk.Contains([]byte{
		INTEGER8,
		INTEGER16,
		INTEGER32,
		INTEGER64,
	}, t)
}

func IsUnsignedType(t byte) bool {
	return funk.Contains([]byte{
		UNSIGNED8,
		UNSIGNED16,
		UNSIGNED32,
		UNSIGNED64,
	}, t)
}

func IsIntegerType(t byte) bool {
	return funk.Contains([]byte{
		UNSIGNED8,
		UNSIGNED16,
		UNSIGNED32,
		UNSIGNED64,
		INTEGER8,
		INTEGER16,
		INTEGER32,
		INTEGER64,
	}, t)
}

func IsFloatType(t byte) bool {
	return funk.Contains([]byte{
		REAL32,
		REAL64,
	}, t)
}

func IsNumberType(t byte) bool {
	return funk.Contains([]byte{
		UNSIGNED8,
		UNSIGNED16,
		UNSIGNED32,
		UNSIGNED64,
		INTEGER8,
		INTEGER16,
		INTEGER32,
		INTEGER64,
		REAL32,
		REAL64,
	}, t)
}

func IsDataType(t byte) bool {
	return funk.Contains([]byte{
		VISIBLE_STRING,
		OCTET_STRING,
		UNICODE_STRING,
		DOMAIN,
	}, t)
}
