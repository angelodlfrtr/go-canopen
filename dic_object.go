package canopen

type DicObject interface {
	// For DicRecord and DicArray

	GetIndex() uint16
	GetSubIndex() uint8
	GetName() string
	AddMember(DicObject)
	FindIndex(uint16) DicObject
	FindName(string) DicObject

	// For DicVariable only

	GetDataType() byte
	GetDataLen() int

	SetSize(int)
	SetOffset(int)

	Read() error
	Save() error

	GetData() []byte
	SetData([]byte)

	GetStringVal() *string
	GetFloatVal() *float64
	GetUintVal() *uint64
	GetIntVal() *int64
	GetBoolVal() *bool
	GetByteVal() *byte

	SetStringVal(string)
	SetFloatVal(float64)
	SetUintVal(uint64)
	SetIntVal(int64)
	SetBoolVal(bool)
	SetByteVal(byte)

	// For All DicObject

	IsDicVariable() bool
	SetSDO(*SDOClient)
}
