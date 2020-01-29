package canopen

// Node is a canopen node
type Node struct {
	// Each node has an id, which is ArbitrationID & 0x7F
	ID int

	Network   *Network
	ObjectDic *DicObjectDic

	SDOClient *SDOClient
	PDONode   *PDONode
	NMTMaster *NMTMaster
}

func NewNode(id int, network *Network, objectDic *DicObjectDic) *Node {
	node := &Node{
		ID:        id,
		Network:   network,
		ObjectDic: objectDic,
	}

	return node
}

// SetNetwork set node.Network to the desired network
func (node *Node) SetNetwork(network *Network) {
	node.Network = network
}

// SetObjectDic set node.ObjectDic to the desired ObjectDic
func (node *Node) SetObjectDic(objectDic *DicObjectDic) {
	node.ObjectDic = objectDic
}

// Init create sdo clients, pdo nodes, nmt master
func (node *Node) Init() {
	node.SDOClient = NewSDOClient(node)
	node.PDONode = NewPDONode(node)
	node.NMTMaster = NewNMTMaster(node.ID, node.Network)

	// @TODO: list for NMTMaster
	// @TODO: implement EMCY
}

// Stop node
func (node *Node) Stop() {
	// Stop nmt master
	node.NMTMaster.UnlistenForHeartbeat()

	// Stop pdo listeners
	for _, mm := range node.PDONode.RX.Maps {
		mm.Unlisten()
	}
	for _, mm := range node.PDONode.TX.Maps {
		mm.Unlisten()
	}
}
