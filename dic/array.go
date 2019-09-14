package dic

import (
	"fmt"
)

type Array struct {
	Description string
	Index       int
	Name        string

	SubIndexes map[int]Object
	SubNames   map[string]int
}

func (array *Array) GetIndex() int {
	return array.Index
}

func (array *Array) GetName() string {
	return array.Name
}

func (array *Array) AddMember(object Object) {
	array.SubIndexes[object.GetIndex()] = object
	array.SubNames[object.GetName()] = object.GetIndex()
}

func (array *Array) FindIndex(index int) (Object, error) {
	if object, ok := array.SubIndexes[index]; ok {
		return object, nil
	}

	return nil, fmt.Errorf("Object with index %d not found in array", index)
}

func (array *Array) FindName(name string) (Object, error) {
	if index, ok := array.SubNames[name]; ok {
		return array.SubIndexes[index], nil
	}

	return nil, fmt.Errorf("Object with name %s not found in array", name)
}
