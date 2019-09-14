package dic

import (
	"fmt"
)

type ObjectDic struct {
	Baudrate int
	NodeID   int

	// Map of Object ids to objects
	Indexes map[int]Object

	// Index to map objects names to objects indexs
	NamesIndex map[string]int
}

func (objectDic *ObjectDic) AddObject(object Object) {
	objectDic.Indexes[object.GetIndex()] = object
	objectDic.NamesIndex[object.GetName()] = object.GetIndex()
}

func (objectDic *ObjectDic) FindIndex(index int) (Object, error) {
	if object, ok := objectDic.Indexes[index]; ok {
		return object, nil
	}

	return nil, fmt.Errorf("Object with index %d not found in dictionnary", index)
}

func (objectDic *ObjectDic) FindName(name string) (Object, error) {
	if index, ok := objectDic.NamesIndex[name]; ok {
		return objectDic.Indexes[index], nil
	}

	return nil, fmt.Errorf("Object with name %s not found in dictionnary", name)
}
