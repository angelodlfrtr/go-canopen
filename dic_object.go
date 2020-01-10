package canopen

type DicObject interface {
	GetIndex() uint16
	GetSubIndex() uint8
	GetName() string
	AddMember(DicObject)
	FindIndex(uint16) DicObject
	FindName(string) DicObject

	SetSDO(*SDOClient)

	IsDicVariable() bool

	GetDataLen() int
	SetSize(int)
	SetOffset(int)

	Read() error
	GetData() []byte
	GetStringVal() *string
	GetFloatVal() *float64
	GetUintVal() *uint64
	GetIntVal() *int64
	GetBoolVal() *bool
	GetByteVal() *byte
}
