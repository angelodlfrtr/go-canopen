package canopen

import (
	"encoding/ascii85"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

type DicVariable struct {
	Unit        string
	Factor      int
	Min         int
	Max         int
	Default     []byte
	DataType    byte
	AccessType  string
	Description string

	SDOClient *SDOClient

	Data   []byte
	Offset int
	Size   int

	Index    uint16
	SubIndex uint8
	Name     string

	ValueDescriptions map[string]string
	BitDefinitions    map[string][]byte
}

func (variable *DicVariable) GetIndex() uint16 {
	return variable.Index
}

func (variable *DicVariable) GetSubIndex() uint8 {
	return variable.SubIndex
}

func (variable *DicVariable) GetName() string {
	return variable.Name
}

func (variable *DicVariable) AddMember(object DicObject) {
	// Do nothing here
}

func (variable *DicVariable) FindIndex(index uint16) DicObject {
	// Do nothing here
	return nil
}

func (variable *DicVariable) FindName(name string) DicObject {
	// Do nothing here
	return nil
}

func (variable *DicVariable) SetSDO(sdo *SDOClient) {
	variable.SDOClient = sdo
}

func (variable *DicVariable) IsDicVariable() bool {
	return true
}

func (variable *DicVariable) GetDataType() byte {
	return variable.DataType
}

func (variable *DicVariable) GetDataLen() int {
	l := 1

	if variable.DataType == Boolean {
		l = 1
	}

	if variable.DataType == Integer8 {
		l = 1
	}

	if variable.DataType == Integer16 {
		l = 2
	}

	if variable.DataType == Integer32 {
		l = 4
	}

	if variable.DataType == Integer64 {
		l = 8
	}

	if variable.DataType == Unsigned8 {
		l = 1
	}

	if variable.DataType == Unsigned16 {
		l = 2
	}

	if variable.DataType == Unsigned32 {
		l = 4
	}

	if variable.DataType == Unsigned64 {
		l = 8
	}

	if variable.DataType == Real32 {
		l = 4
	}

	if variable.DataType == Real64 {
		l = 8
	}

	return l * 8
}

func (variable *DicVariable) SetSize(s int) {
	variable.Size = s
}

func (variable *DicVariable) SetOffset(s int) {
	variable.Offset = s
}

func (variable *DicVariable) GetOffset() int {
	return variable.Offset
}

func (variable *DicVariable) AddValueDescription(name string, des string) {
	variable.ValueDescriptions[name] = des
}

func (variable *DicVariable) AddBitDefinition(name string, bits []byte) {
	variable.BitDefinitions[name] = bits
}

// Read variable value using SDO
func (variable *DicVariable) Read() error {
	if variable.SDOClient == nil {
		return errors.New("SDOClient required")
	}

	data, err := variable.SDOClient.Read(variable.Index, variable.SubIndex)
	if err != nil {
		return err
	}

	variable.Data = data

	return nil
}

// Write variable value using SDO
func (variable *DicVariable) Write(data []byte) error {
	if variable.SDOClient == nil {
		return errors.New("SDOClient required")
	}

	return variable.SDOClient.Write(
		variable.Index,
		variable.SubIndex,
		variable.IsDomainDataType(),
		variable.Data,
	)
}

// Save variable.Data using SDO
func (variable *DicVariable) Save() error {
	return variable.Write(variable.Data)
}

func (variable *DicVariable) IsDomainDataType() bool {
	return variable.DataType == Domain
}

func (variable *DicVariable) GetData() []byte {
	return variable.Data
}

func (variable *DicVariable) SetData(data []byte) {
	variable.Data = data
}

func (variable *DicVariable) GetStringVal() *string {
	if !IsStringType(variable.DataType) {
		return nil
	}

	// @TODO: check if all working
	var v string

	if variable.DataType == VisibleString {
		dst := []byte{}
		ascii85.Decode(dst, variable.Data, false)
		v = string(dst)
	}

	if variable.DataType == UnicodeString {
		src := make([]uint16, len(variable.Data))
		for i, r := range variable.Data {
			src[i] = uint16(r)
		}
		dst := utf16.Decode(src)
		v = string(dst)
	}

	if len(v) == 0 {
		v = string(variable.Data)
	}

	return &v
}

func (variable *DicVariable) GetFloatVal() *float64 {
	if !IsFloatType(variable.DataType) {
		return nil
	}

	var v float64

	// @TODO

	return &v
}

func (variable *DicVariable) GetUintVal() *uint64 {
	if !IsUnsignedType(variable.DataType) {
		return nil
	}

	var v uint64

	if variable.DataType == Unsigned8 {
		v = uint64(variable.Data[0])
	}

	if variable.DataType == Unsigned16 {
		v = uint64(binary.LittleEndian.Uint16(variable.Data))
	}

	if variable.DataType == Unsigned32 {
		v = uint64(binary.LittleEndian.Uint32(variable.Data))
	}

	if variable.DataType == Unsigned64 {
		v = binary.LittleEndian.Uint64(variable.Data)
	}

	return &v
}

func (variable *DicVariable) GetIntVal() *int64 {
	if !IsSignedType(variable.DataType) {
		return nil
	}

	var v int64

	if variable.DataType == Integer8 {
		v = int64(int8(variable.Data[0]))
	}

	// @TODO: https://groups.google.com/forum/#!topic/golang-nuts/f1QQkP19G9Q
	if variable.DataType == Integer16 {
		v = int64(binary.LittleEndian.Uint16(variable.Data))
	}

	if variable.DataType == Integer32 {
		v = int64(binary.LittleEndian.Uint32(variable.Data))
	}

	if variable.DataType == Integer64 {
		v = int64(binary.LittleEndian.Uint64(variable.Data))
	}

	return &v
}

func (variable *DicVariable) GetBoolVal() *bool {
	if variable.DataType != Boolean {
		return nil
	}

	v := false

	// @TODO

	return &v
}

func (variable *DicVariable) GetByteVal() *byte {
	var v byte

	if variable.DataType == Unsigned8 {
		v = variable.Data[0]
	}

	return &v
}

func (variable *DicVariable) SetStringVal(a string) {
	// @TODO
}

func (variable *DicVariable) SetFloatVal(a float64) {
	// @TODO
}

func (variable *DicVariable) SetUintVal(a uint64) {
	if !IsUnsignedType(variable.DataType) {
		return
	}

	if variable.DataType == Unsigned8 {
		variable.Data[0] = byte(a)
	}

	if variable.DataType == Unsigned16 {
		binary.LittleEndian.PutUint16(variable.Data, uint16(a))
	}

	if variable.DataType == Unsigned32 {
		binary.LittleEndian.PutUint32(variable.Data, uint32(a))
	}

	if variable.DataType == Unsigned64 {
		binary.LittleEndian.PutUint64(variable.Data, uint64(a))
	}
}

func (variable *DicVariable) SetIntVal(a int64) {
	if !IsSignedType(variable.DataType) {
		return
	}

	if variable.DataType == Integer8 {
		variable.Data[0] = byte(a)
	}

	// @TODO: https://groups.google.com/forum/#!topic/golang-nuts/f1QQkP19G9Q
	if variable.DataType == Integer16 {
		binary.LittleEndian.PutUint16(variable.Data, uint16(a))
	}

	if variable.DataType == Integer32 {
		binary.LittleEndian.PutUint32(variable.Data, uint32(a))
	}

	if variable.DataType == Integer64 {
		binary.LittleEndian.PutUint64(variable.Data, uint64(a))
	}
}

func (variable *DicVariable) SetBoolVal(a bool) {
	if variable.DataType == Boolean {
		// @TODO
	}
}

func (variable *DicVariable) SetByteVal(a byte) {
	if variable.DataType == Unsigned8 {
		variable.Data[0] = a
	}
}
