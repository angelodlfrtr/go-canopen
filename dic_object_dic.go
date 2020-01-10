package canopen

type DicObjectDic struct {
	Baudrate int
	NodeID   int

	// Map of Object ids to objects
	Indexes map[uint16]DicObject

	// Index to map objects names to objects indexs
	NamesIndex map[string]uint16
}

func NewDicObjectDic() *DicObjectDic {
	return &DicObjectDic{
		Indexes:    map[uint16]DicObject{},
		NamesIndex: map[string]uint16{},
	}
}

func (objectDic *DicObjectDic) AddObject(object DicObject) {
	objectDic.Indexes[object.GetIndex()] = object
	objectDic.NamesIndex[object.GetName()] = object.GetIndex()
}

func (objectDic *DicObjectDic) FindIndex(index uint16) DicObject {
	if object, ok := objectDic.Indexes[index]; ok {
		return object
	}

	return nil
}

func (objectDic *DicObjectDic) FindName(name string) DicObject {
	if index, ok := objectDic.NamesIndex[name]; ok {
		return objectDic.Indexes[index]
	}

	return nil
}
