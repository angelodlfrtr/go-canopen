package canopen

// Maps define a PDO map
type PDOMaps struct {
	PDONode *PDONode
	Maps    map[int]*PDOMap
}

// NewMap create a new Maps
func NewPDOMaps(comOffset, mapOffset int, pdoNode *PDONode) *PDOMaps {
	pdoMaps := &PDOMaps{
		PDONode: pdoNode,
		Maps:    make(map[int]*PDOMap),
	}

	for i := 0; i < 32; i++ {
		if comSdo := pdoMaps.PDONode.Node.ObjectDic.FindIndex(uint16(comOffset + i)); comSdo != nil {
			mapSdo := pdoMaps.PDONode.Node.ObjectDic.FindIndex(uint16(mapOffset + i))

			comSdo.SetSDO(pdoMaps.PDONode.Node.SDOClient)
			mapSdo.SetSDO(pdoMaps.PDONode.Node.SDOClient)

			pdoMaps.Maps[i+1] = NewPDOMap(pdoNode, comSdo, mapSdo)
		}
	}

	return pdoMaps
}

// Find a map by index
func (maps *PDOMaps) FindIndex(idx int) *PDOMap {
	if m, ok := maps.Maps[idx]; ok {
		return m
	}

	return nil
}

// FindByName a map
func (maps *PDOMaps) FindByName(name string) *PDOMap {
	var m *PDOMap

	for _, ma := range maps.Maps {
		for _, v := range ma.Map {
			if v.GetName() == name {
				m = ma
				break
			}
		}

		if m != nil {
			break
		}
	}

	return m
}
