package canopen

type DicRecord struct {
	Description string
	Index       uint16
	Name        string

	SDOClient *SDOClient

	SubIndexes map[uint8]DicObject
	SubNames   map[string]uint8
}

// GetIndex of DicRecord
func (record *DicRecord) GetIndex() uint16 {
	return record.Index
}

// GetSubIndex not applicable
func (record *DicRecord) GetSubIndex() uint8 { return 0 }

// GetName of DicRecord
func (record *DicRecord) GetName() string {
	return record.Name
}

// AddMember to DicRecord
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

// FindIndex find by index a DicObject in DicRecord
func (record *DicRecord) FindIndex(index uint16) DicObject {
	if object, ok := record.SubIndexes[uint8(index)]; ok {
		object.SetSDO(record.SDOClient)
		return object
	}

	return nil
}

// FindName find by name a DicObject in DicRecord
func (record *DicRecord) FindName(name string) DicObject {
	if index, ok := record.SubNames[name]; ok {
		return record.FindIndex(uint16(index))
	}

	return nil
}

func (record *DicRecord) GetDataType() byte     { return 0x00 }
func (record *DicRecord) GetDataLen() int       { return 0 }
func (record *DicRecord) SetSize(s int)         {}
func (record *DicRecord) SetOffset(s int)       {}
func (record *DicRecord) Read() error           { return nil }
func (record *DicRecord) Save() error           { return nil }
func (record *DicRecord) GetData() []byte       { return nil }
func (record *DicRecord) SetData(data []byte)   {}
func (record *DicRecord) GetStringVal() *string { return nil }
func (record *DicRecord) GetFloatVal() *float64 { return nil }
func (record *DicRecord) GetUintVal() *uint64   { return nil }
func (record *DicRecord) GetIntVal() *int64     { return nil }
func (record *DicRecord) GetBoolVal() *bool     { return nil }
func (record *DicRecord) GetByteVal() *byte     { return nil }
func (record *DicRecord) SetStringVal(a string) {}
func (record *DicRecord) SetFloatVal(a float64) {}
func (record *DicRecord) SetUintVal(a uint64)   {}
func (record *DicRecord) SetIntVal(a int64)     {}
func (record *DicRecord) SetBoolVal(a bool)     {}
func (record *DicRecord) SetByteVal(a byte)     {}
func (record *DicRecord) IsDicVariable() bool   { return false }

// SetSDO to DicRecord
func (record *DicRecord) SetSDO(sdo *SDOClient) {
	record.SDOClient = sdo
}
