package canopen

type PDONode struct {
	Node    *Node
	Network *Network
	RX      *PDOMaps
	TX      *PDOMaps
}

func NewPDONode(n *Node) *PDONode {
	pdoNode := &PDONode{Node: n}
	pdoNode.RX = NewPDOMaps(0x1400, 0x1600, pdoNode)
	pdoNode.TX = NewPDOMaps(0x1800, 0x1A00, pdoNode)

	return pdoNode
}

func (node *PDONode) FindByName(name string) *PDOMap {
	r := node.RX.FindByName(name)
	if r == nil {
		r = node.TX.FindByName(name)
	}

	return r
}

func (node *PDONode) Read() error {
	for _, maps := range []*PDOMaps{node.RX, node.TX} {
		for _, v := range maps.Maps {
			if err := v.Read(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (node *PDONode) Save() error {
	for _, maps := range []*PDOMaps{node.RX, node.TX} {
		for _, v := range maps.Maps {
			if err := v.Save(); err != nil {
				return err
			}
		}
	}

	return nil
}
