package pdo

type Maps struct {
	ComOffset int
	MapOffset int
	Node      *Node
	Maps      map[int]*Map
}

func NewMaps(comOffset, mapOffset int, pdoNode *Node) *Maps {
	pdoMaps := &Maps{
		ComOffset: comOffset,
		MapOffset: mapOffset,
		Node:      pdoNode,
	}

	for i := 0; i < 32; i++ {
		idx := comOffset + i

		if ob, err := pdoNode.Node.ObjectDic.FindIndex(idx); err == nil {
			pdoMaps.Maps[i+1] = NewMap(pdoNode, ob, ob)
		}
	}

	return pdoMaps
}
