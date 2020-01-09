package dic

import (
	"fmt"
)

type Record struct {
	Description string
	Index       int
	Name        string

	SubIndexes map[int]Object
	SubNames   map[string]int
}

func (record *Record) GetIndex() int {
	return record.Index
}

func (record *Record) GetName() string {
	return record.Name
}

func (record *Record) AddMember(object Object) {
	record.SubIndexes[object.GetIndex()] = object
	record.SubNames[object.GetName()] = object.GetIndex()
}

func (record *Record) FindIndex(index int) (Object, error) {
	if object, ok := record.SubIndexes[index]; ok {
		return object, nil
	}

	return nil, fmt.Errorf("object with index %d not found in record", index)
}

func (record *Record) FindName(name string) (Object, error) {
	if index, ok := record.SubNames[name]; ok {
		return record.SubIndexes[index], nil
	}

	return nil, fmt.Errorf("object with name %s not found in record", name)
}
