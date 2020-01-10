package canopen

type DicRecord struct {
	Description string
	Index       uint16
	Name        string

	SubIndexes map[uint8]DicObject
	SubNames   map[string]uint8
}

func (record *DicRecord) GetIndex() uint16 {
	return record.Index
}

func (record *DicRecord) GetSubIndex() uint8 { return 0 }

func (record *DicRecord) GetName() string {
	return record.Name
}

func (record *DicRecord) AddMember(object DicObject) {
	if record.SubIndexes == nil {
		record.SubIndexes = map[uint8]DicObject{}
	}

	if record.SubNames == nil {
		record.SubNames = map[string]uint8{}
	}

	record.SubIndexes[object.GetSubIndex()] = object
	record.SubNames[object.GetName()] = object.GetSubIndex()
}

func (record *DicRecord) FindIndex(index uint16) DicObject {
	if object, ok := record.SubIndexes[uint8(index)]; ok {
		return object
	}

	return nil
}

func (record *DicRecord) FindName(name string) DicObject {
	if index, ok := record.SubNames[name]; ok {
		return record.SubIndexes[index]
	}

	return nil
}

func (record *DicRecord) SetSDO(sdo *SDOClient) {}
func (record *DicRecord) IsDicVariable() bool   { return false }
func (record *DicRecord) GetDataLen() int       { return 0 }
func (record *DicRecord) SetSize(s int)         {}
func (record *DicRecord) SetOffset(s int)       {}
func (record *DicRecord) Read() error           { return nil }
func (record *DicRecord) GetData() []byte       { return nil }
func (record *DicRecord) GetStringVal() *string { return nil }
func (record *DicRecord) GetFloatVal() *float64 { return nil }
func (record *DicRecord) GetUintVal() *uint64   { return nil }
func (record *DicRecord) GetIntVal() *int64     { return nil }
func (record *DicRecord) GetBoolVal() *bool     { return nil }
func (record *DicRecord) GetByteVal() *byte     { return nil }
