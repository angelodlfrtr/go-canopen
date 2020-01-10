package canopen

type DicArray struct {
	Description string
	Index       uint16
	Name        string

	SubIndexes map[uint8]DicObject
	SubNames   map[string]uint8
}

func (array *DicArray) GetIndex() uint16 {
	return array.Index
}

func (array *DicArray) GetSubIndex() uint8 { return 0 }

func (array *DicArray) GetName() string {
	return array.Name
}

// AddMember to SubIndexes
func (array *DicArray) AddMember(object DicObject) {
	if array.SubIndexes == nil {
		array.SubIndexes = map[uint8]DicObject{}
	}

	if array.SubNames == nil {
		array.SubNames = map[string]uint8{}
	}

	array.SubIndexes[object.GetSubIndex()] = object
	array.SubNames[object.GetName()] = object.GetSubIndex()
}

func (array *DicArray) FindIndex(index uint16) DicObject {
	if object, ok := array.SubIndexes[uint8(index)]; ok {
		return object
	}

	return nil
}

func (array *DicArray) FindName(name string) DicObject {
	if index, ok := array.SubNames[name]; ok {
		return array.SubIndexes[index]
	}

	return nil
}

func (array *DicArray) SetSDO(sdo *SDOClient) {}
func (array *DicArray) IsDicVariable() bool   { return false }
func (array *DicArray) GetDataLen() int       { return 0 }
func (array *DicArray) SetSize(s int)         {}
func (array *DicArray) SetOffset(s int)       {}
func (array *DicArray) Read() error           { return nil }
func (array *DicArray) GetData() []byte       { return nil }
func (array *DicArray) GetStringVal() *string { return nil }
func (array *DicArray) GetFloatVal() *float64 { return nil }
func (array *DicArray) GetUintVal() *uint64   { return nil }
func (array *DicArray) GetIntVal() *int64     { return nil }
func (array *DicArray) GetBoolVal() *bool     { return nil }
func (array *DicArray) GetByteVal() *byte     { return nil }